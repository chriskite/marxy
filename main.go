package main

import (
	"time"
)

func main() {
	run()
}

func run() {
	host := "http://localhost"
	marathon := NewMarathon(host)
	configChan := make(chan string, 10)
	runChan := make(chan bool)
	go configActor(configChan)
	go marathonActor(marathon, runChan, configChan)
	timingActor(runChan)
}

func timingActor(runChan chan<- bool) {
	for {
		runChan <- true
		time.Sleep(time.Second * 5)
	}
}
