package main

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
