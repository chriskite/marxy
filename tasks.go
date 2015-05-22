package main

import (
	"fmt"
	"strings"
)

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	AppId              string              `json:"appId"`
	Id                 string              `json:"id"`
	Host               string              `json:"host"`
	Ports              []int64             `json:"ports"`
	StartedAt          string              `json:"startedAt"`
	StagedAt           string              `json:"stagedAt"`
	Version            string              `json:"version"`
	ServicePorts       []int64             `json:"servicePorts"`
	HealthCheckResults []HealthCheckResult `json:"healthCheckResults"`
}

type HealthCheckResult struct {
	TaskId              string `json:"taskId"`
	FirstSuccess        string `json:"firstSuccess"`
	LastSuccess         string `json:"lastSuccess"`
	LastFailure         string `json:"lastFailure"`
	ConsecutiveFailures int64  `json:"consecutiveFailures"`
	Alive               bool   `json:"alive"`
}

func (t Task) EscapedAppId() string {
	return strings.Replace(strings.TrimLeft(t.AppId, "/"), "/", "_", -1)
}

func (t Task) ServerLine(portIndex, serverIndex int) (string, error) {
	if portIndex < 0 || portIndex >= len(t.Ports) {
		return "", fmt.Errorf("portIndex %d out of range", portIndex)
	}

	port := t.Ports[portIndex]
	return fmt.Sprintf("server %s-%d %s:%d check maxconn 0", t.EscapedAppId(), serverIndex, t.Host, port), nil
}

func (t Task) IsAlive() bool {
	hcr := t.HealthCheckResults
	if hcr == nil || len(hcr) == 0 {
		return true
	}

	return hcr[0].Alive
}
