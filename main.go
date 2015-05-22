package main

import (
	"bytes"
	"fmt"
)

func main() {
	host := "http://marathon.ocean"
	marathon := NewMarathon(host)
	config, err := getHAProxyConfig(marathon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)
}

func getHAProxyConfig(marathon Marathon) (string, error) {
	tasksResp, err := GetTasks(marathon)
	if err != nil {
		return "", err
	}
	taskMap := make(map[string][]Task)
	for _, task := range tasksResp.Tasks {
		slice, ok := taskMap[task.EscapedAppId()]
		if !ok {
			slice = make([]Task, 0)
		}
		taskMap[task.EscapedAppId()] = append(slice, task)
	}

	var buffer bytes.Buffer
	// TODO put haproxy header in buffer

	// TODO this is wrong in the sense that it doesnt account for multiple service ports
	for appId, tasks := range taskMap {
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
			// put service header in buffer
			servicePort := tasks[0].ServicePorts[0]
			buffer.WriteString(fmt.Sprintf("listen %s-%d\n", appId, servicePort))
			buffer.WriteString(fmt.Sprintf("  bind 0.0.0.0:%d\n", servicePort))
			buffer.WriteString("  mode tcp\n  option tcplog\n  balance leastconn\n")

			// put each server line in buffer
			for _, line := range lines {
				buffer.WriteString("  " + line + "\n")
			}
			buffer.WriteString("\n")
		}
	}

	// output buffer
	return buffer.String(), nil
}
