package main

import (
	"bytes"
	"fmt"
	"sort"
)

func configActor(configChan <-chan string) {
	var oldConfig string
	for {
		select {
		case config := <-configChan:
			handleConfig(config, oldConfig)
			oldConfig = config
		}
	}
}

func handleConfig(config string, oldConfig string) {
	if config != oldConfig {
		fmt.Println(config)
		// write temp file
		// atomic mv temp file to config file
		// reload haproxy
	}
}

func haproxyConfigHeader() string {
	header := `global
  daemon
  maxconn 4096

defaults
  log                 global
  retries             3
  maxconn             1024
  timeout connect     5s
  timeout client      60s
  timeout server      60s
  timeout client-fin  60s
  timeout tunnel      12h

listen stats :9090
    mode http
    stats enable
    stats realm HAProxy\ Statistics
    stats uri /

`
	return header
}

func haproxyConfig(marathonTasks []Task) (string, error) {
	sort.Sort(TaskSlice(marathonTasks))

	// make a map from appId to slice of tasks in that app
	appMap := make(map[string][]Task)
	for _, task := range marathonTasks {
		slice, ok := appMap[task.EscapedAppId()]
		if !ok {
			slice = make([]Task, 0)
		}
		appMap[task.EscapedAppId()] = append(slice, task)
	}

	// buffer containing the haproxy config
	var config bytes.Buffer

	config.WriteString(haproxyConfigHeader())

	// sort the keys so the config file ordering is deterministic
	keys := make([]string, len(appMap))
	i := 0
	for k := range appMap {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)

	for _, appId := range keys {
		tasks := appMap[appId]
		if tasks == nil || len(tasks) == 0 {
			continue
		}

		// foreach service port in the first task
		for portIndex, servicePort := range tasks[0].ServicePorts {
			lines := make([]string, 0)

			for serverIndex, task := range tasks {
				if !task.IsAlive() {
					continue
				}
				line, err := task.ServerLine(portIndex, serverIndex)
				if err != nil {
					continue
				}
				lines = append(lines, line)
			}

			if len(lines) > 0 {
				// put service header in config
				config.WriteString(fmt.Sprintf("listen %s-%d\n", appId, servicePort))
				config.WriteString(fmt.Sprintf("  bind 0.0.0.0:%d\n", servicePort))
				config.WriteString("  mode tcp\n  option tcplog\n  balance leastconn\n")

				// put each server line in config
				for _, line := range lines {
					config.WriteString("  " + line + "\n")
				}
				config.WriteString("\n")
			}
		}
	}

	// output config
	return config.String(), nil
}
