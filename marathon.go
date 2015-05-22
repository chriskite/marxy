package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Marathon interface {
	getTasksReq() (*http.Request, error)
}

type UnauthMarathon struct {
	baseUrl string
}

type BasicAuthMarathon struct {
	baseUrl  string
	username string
	password string
}

func NewMarathon(url string) *UnauthMarathon {
	return &UnauthMarathon{baseUrl: url}
}
func NewAuthMarathon(url, user, pass string) *BasicAuthMarathon {
	return &BasicAuthMarathon{baseUrl: url, username: user, password: pass}
}

func unauthTasksReq(baseUrl string) (*http.Request, error) {
	url := fmt.Sprintf("%s/v2/tasks", baseUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return req, err
	}

	return req, nil
}

func (m *BasicAuthMarathon) getTasksReq() (*http.Request, error) {
	req, err := unauthTasksReq(m.baseUrl)
	if err != nil {
		return req, err
	}

	req.SetBasicAuth(m.username, m.password)

	return req, nil
}

func (m *UnauthMarathon) getTasksReq() (*http.Request, error) {
	return unauthTasksReq(m.baseUrl)
}

func GetTasks(m Marathon) (TasksResponse, error) {
	var tasks TasksResponse

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := m.getTasksReq()
	if err != nil {
		return tasks, err
	}

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return tasks, fmt.Errorf("HTTP request failed (%d)", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tasks, err
	}

	err = json.Unmarshal(body, &tasks)
	return tasks, err
}