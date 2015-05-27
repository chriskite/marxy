package main

import (
	"github.com/codegangsta/cli"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "marxy"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "marathon",
			Usage:  "marathon url e.g. http://marathon",
			EnvVar: "MARXY_MARATHON",
		},
		cli.StringFlag{
			Name:   "http-user",
			Usage:  "http basic auth user",
			EnvVar: "MARXY_HTTP_USER",
		},
		cli.StringFlag{
			Name:   "http-pass",
			Usage:  "http basic auth pass",
			EnvVar: "MARXY_HTTP_PASS",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	host := c.String("marathon")
	if host == "" {
		return
	}

	var marathon Marathon

	user, pass := c.String("http-user"), c.String("http-pass")

	if user != "" && pass != "" {
		marathon = NewAuthMarathon(host, user, pass)
	} else {
		marathon = NewMarathon(host)
	}

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
