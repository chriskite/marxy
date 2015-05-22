package main

import (
	"bytes"
	"fmt"
)

func main() {
	host := "http://marathon.ocean"
	marathon := NewMarathon(host)
	config, err := haproxyConfig(marathon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)
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

func haproxyConfig(marathon Marathon) (string, error) {
	tasksResp, err := GetTasks(marathon)
	if err != nil {
		return "", err
	}
	// make a map from appId to slice of tasks in that app
	appMap := make(map[string][]Task)
	for _, task := range tasksResp.Tasks {
		slice, ok := appMap[task.EscapedAppId()]
		if !ok {
			slice = make([]Task, 0)
		}
		appMap[task.EscapedAppId()] = append(slice, task)
	}

	// buffer containing the haproxy config
	var config bytes.Buffer

	config.WriteString(haproxyConfigHeader())

	// FIXME this is wrong in the sense that it doesnt account for multiple service ports
	for appId, tasks := range appMap {
		lines := make([]string, 0)

		i := 0
		for _, task := range tasks {
			if !task.IsAlive() {
				continue
			}
			line, err := task.ServerLine(0, i)
			if err != nil {
				continue
			}
			lines = append(lines, line)
		}

		if len(lines) > 0 {
			// put service header in config
			servicePort := tasks[0].ServicePorts[0]
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

	// output config
	return config.String(), nil
}
