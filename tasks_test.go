package main

import (
	"encoding/json"
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TasksS struct{}

var _ = Suite(&TasksS{})

func testJson() []byte {
	json := []byte(`{
	"tasks": [
		{
			"appId": "/task1",
			"id": "task1",
			"host": "ip-10-0-0-2.ec2.internal",
			"ports": [31004,31005],
			"startedAt": "2015-05-06T18:34:31.753Z",
			"stagedAt": "2015-05-06T18:34:31.354Z",
			"version": "2015-04-29T21:03:19.522Z",
			"servicePorts": [12002,10003],
			"healthCheckResults": [
				{
					"taskId": "task1",
					"firstSuccess": "2015-05-06T18:34:45.386Z",
					"lastSuccess": "2015-05-21T18:19:03.705Z",
					"lastFailure": null,
					"consecutiveFailures": 0,
					"alive": true
				}
			]
		},
		{
			"appId": "/task2",
			"id": "task2",
			"host": "ip-10-0-0-3.ec2.internal",
			"ports": [31180],
			"startedAt": "2015-05-20T22:12:17.019Z",
			"stagedAt": "2015-05-20T22:10:22.376Z",
			"version": "2015-05-20T22:10:16.322Z",
			"servicePorts": [10000]
		}
	]
}`)
	return json
}

func simpleTask() Task {
	task := Task{
		AppId:              "foo",
		Id:                 "foo.1",
		Host:               "localhost",
		Ports:              []int64{31000},
		StartedAt:          "2015-05-20T22:12:17.019Z",
		StagedAt:           "2015-05-20T22:10:22.376Z",
		Version:            "2015-05-20T22:10:16.322Z",
		ServicePorts:       []int64{10000},
		HealthCheckResults: nil,
	}
	return task
}

func (s *TasksS) TestUnmarshalTasks(c *C) {
	var r TasksResponse
	err := json.Unmarshal(testJson(), &r)
	c.Assert(err, IsNil)

	tasks := r.Tasks
	c.Assert(tasks, HasLen, 2)

	c.Check(tasks[0], DeepEquals, Task{
		AppId:        "/task1",
		Id:           "task1",
		Host:         "ip-10-0-0-2.ec2.internal",
		Ports:        []int64{31004, 31005},
		StartedAt:    "2015-05-06T18:34:31.753Z",
		StagedAt:     "2015-05-06T18:34:31.354Z",
		Version:      "2015-04-29T21:03:19.522Z",
		ServicePorts: []int64{12002, 10003},
		HealthCheckResults: []HealthCheckResult{
			{
				TaskId:              "task1",
				FirstSuccess:        "2015-05-06T18:34:45.386Z",
				LastSuccess:         "2015-05-21T18:19:03.705Z",
				LastFailure:         "",
				ConsecutiveFailures: 0,
				Alive:               true,
			},
		},
	})

	c.Check(tasks[1], DeepEquals, Task{
		AppId:              "/task2",
		Id:                 "task2",
		Host:               "ip-10-0-0-3.ec2.internal",
		Ports:              []int64{31180},
		StartedAt:          "2015-05-20T22:12:17.019Z",
		StagedAt:           "2015-05-20T22:10:22.376Z",
		Version:            "2015-05-20T22:10:16.322Z",
		ServicePorts:       []int64{10000},
		HealthCheckResults: nil,
	})
}

func (s *TasksS) TestServerLine(c *C) {
	task := simpleTask()
	expectedLine := "server foo-10000-0 localhost:31000 check maxconn 0"
	line, err := task.ServerLine(0, 0)
	c.Assert(err, IsNil)
	c.Check(line, Equals, expectedLine)
}

func (s *TasksS) TestServerLineBadIndex(c *C) {
	task := simpleTask()
	_, err := task.ServerLine(1, 0)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "portIndex 1 out of range")
}

func (s *TasksS) TestIsAlive(c *C) {
	// start with empty HealthCheckResults
	task := simpleTask()
	c.Check(task.IsAlive(), Equals, true)

	// set Alive false
	task.HealthCheckResults = []HealthCheckResult{
		{
			Alive: false,
		},
	}
	c.Check(task.IsAlive(), Equals, false)

	// set Alive true
	task.HealthCheckResults = []HealthCheckResult{
		{
			Alive: true,
		},
	}
	c.Check(task.IsAlive(), Equals, true)
}
