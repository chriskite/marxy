package main

import (
	"encoding/json"
	"fmt"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type MarathonS struct{}

var _ = Suite(&MarathonS{})

type MarathonTasksHandler struct{}

func (m MarathonTasksHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/v2/tasks" {
		json := testJson()
		rw.Write(json)
	} else {
		rw.WriteHeader(404)
	}
}

type AuthMarathonTasksHandler struct{}

func (m AuthMarathonTasksHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/v2/tasks" {
		user, pass, ok := r.BasicAuth()
		if user == "guest" && pass == "letmein" && ok {
			json := testJson()
			rw.Write(json)
		} else {
			rw.WriteHeader(403)
		}
	} else {
		rw.WriteHeader(404)
	}
}

type MarathonFailHandler struct{}

func (m MarathonFailHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(500)
	rw.Write([]byte("Internal server error"))
}

func (s *MarathonS) TestGetTasksExample(c *C) {
	server := httptest.NewServer(new(MarathonTasksHandler))
	defer server.Close()

	url := fmt.Sprintf("%s/v2/tasks", server.URL)
	resp, err := http.DefaultClient.Get(url)
	c.Assert(err, IsNil)
	c.Assert(resp, NotNil)
	c.Check(resp.StatusCode, Equals, 200)
	body, err := ioutil.ReadAll(resp.Body)
	c.Check(body, DeepEquals, testJson())
}

func (s *MarathonS) TestGetTasksSuccess(c *C) {
	server := httptest.NewServer(new(MarathonTasksHandler))
	defer server.Close()

	m := NewMarathon(server.URL)
	resp, err := GetTasks(m)
	c.Assert(err, IsNil)
	var expectedTasks TasksResponse
	json.Unmarshal(testJson(), &expectedTasks)
	c.Check(resp, DeepEquals, expectedTasks)
}

func (s *MarathonS) TestGetTasks500(c *C) {
	server := httptest.NewServer(new(MarathonFailHandler))
	defer server.Close()

	m := NewMarathon(server.URL)
	_, err := GetTasks(m)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "HTTP request failed (500)")
}

func (s *MarathonS) TestAuthGetTasksSuccess(c *C) {
	server := httptest.NewServer(new(AuthMarathonTasksHandler))
	defer server.Close()

	m := NewAuthMarathon(server.URL, "guest", "letmein")
	resp, err := GetTasks(m)
	c.Assert(err, IsNil)
	var expectedTasks TasksResponse
	json.Unmarshal(testJson(), &expectedTasks)
	c.Check(resp, DeepEquals, expectedTasks)
}

func (s *MarathonS) TestAuthGetTasks403(c *C) {
	server := httptest.NewServer(new(AuthMarathonTasksHandler))
	defer server.Close()

	m := NewAuthMarathon(server.URL, "foo", "foo")
	_, err := GetTasks(m)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "HTTP request failed (403)")
}
