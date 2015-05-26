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
	if portIndex < 0 || portIndex >= len(t.Ports) || portIndex >= len(t.ServicePorts) {
		return "", fmt.Errorf("portIndex %d out of range", portIndex)
	}

	port := t.Ports[portIndex]
	servicePort := t.ServicePorts[portIndex]
	return fmt.Sprintf(
		"server %s-%d-%d %s:%d check maxconn 0",
		t.EscapedAppId(),
		servicePort,
		serverIndex,
		t.Host,
		port,
	), nil
}

// A Task is alive iff it has passed and is currently passing health checks
func (t Task) IsAlive() bool {
	hcr := t.HealthCheckResults
	if hcr == nil || len(hcr) == 0 {
		return false
	}

	return hcr[0].Alive
}

func (t Task) SortKey() string {
	return t.AppId + t.Id + t.StartedAt + t.Host
}

// Implement sort.Interface
type TaskSlice []Task

func (t TaskSlice) Len() int           { return len(t) }
func (t TaskSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TaskSlice) Less(i, j int) bool { return t[i].SortKey() < t[j].SortKey() }
