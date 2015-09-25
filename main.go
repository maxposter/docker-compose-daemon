package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-compose-daemon"
	app.Usage = "Starting docker containers via docker-compose, redirect docker-compose logs to stderr and stdout, monitoring container state"
	app.Author = "Maxposter"
	app.Email = "development@maxposter.ru"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "configuration, f",
			Usage: "Docker compose config file: -f /path/to/docker-compose.yml",
		},
		cli.StringSliceFlag{
			Name:  "container, c",
			Usage: "Full container name: -c demo_app_1 -c demo_db_1 -c demo_web_1",
		},
		cli.IntFlag{
			Name:  "timeout, t",
			Usage: "Timeout for container monitoring",
			Value: 5,
		},
	}
	app.Action = func(c *cli.Context) {
		var err error

		config := c.String("configuration")
		containers := c.StringSlice("container")
		if config == "" || len(containers) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		startCmd := exec.Command("docker-compose", "-f", config, "up", "-d")
		startCmd.Stdout = os.Stdout
		startCmd.Stderr = os.Stderr

		stopCmd := exec.Command("docker-compose", "-f", config, "stop")
		stopCmd.Stdout = os.Stdout
		stopCmd.Stderr = os.Stderr

		logCmd := exec.Command("docker-compose", "-f", config, "logs")
		logCmd.Stdout = os.Stdout
		logCmd.Stderr = os.Stderr

		err = startCmd.Run()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Cannot start containers: %s\n", err.Error()))
			os.Exit(1)
		}

		go func() {
			err = logCmd.Run()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Cannot read logs: %s\n", err.Error()))
				os.Exit(1)
			}
		}()

		ticker := time.NewTicker(time.Second * time.Duration(c.Int("timeout")))

		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, os.Interrupt)
		signal.Notify(quitChan, syscall.SIGTERM)

		terminating := false
		gracefulStop := func(code int) {
			terminating = true
			os.Stdout.WriteString("\nGraceful stopping containers...\n")
			err = stopCmd.Run()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Cannot graceful stopping containers: %s\n", err.Error()))
				os.Exit(1)
			}
			os.Exit(code)
		}

		for {
			select {
			case <-ticker.C:
				if !terminating {
					for _, container := range containers {
						processCmd := exec.Command("docker", "ps", "-a", "--format", "\"{{.Status}}\"", "-f", "name="+container)
						output, psErr := processCmd.Output()
						if psErr != nil {
							os.Stderr.WriteString(fmt.Sprintf("Cannot check containers: %s\n", psErr.Error()))
							gracefulStop(1)
						}

						status := string(output)
						if !strings.Contains(status, "Up") {
							os.Stderr.WriteString(fmt.Sprintf("Container \"%s\" is down. Current status: %s\n", container, status))
							gracefulStop(1)
						}
					}
				}
			case <-quitChan:
				gracefulStop(0)
			}
		}
	}

	app.Run(os.Args)
}
