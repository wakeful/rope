package main

import (
	"context"
	"flag"
	"os"
	rconfig "rope/pkg/config"
	"rope/pkg/helpers"
	"rope/pkg/logic"
	"time"
)

import (
	"github.com/sirupsen/logrus"
)

var (
	verbose  = flag.Bool("vvv", false, "set log lvl to Debug")
	jsonLogs = flag.Bool("log-json", false, "set log format to json")
	watch    = flag.Bool("watch", false, "watch for config file changes")
	cmdStop  = flag.Bool("stop", false, "stop containers from file")
	cmdList  = flag.Bool("list", false, "list running containers")
)

func main() {
	flag.Parse()

	logrus.SetLevel(logrus.InfoLevel)
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *jsonLogs {
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: false,
			PrettyPrint:      false,
		})
	}

	if 2 > len(os.Args) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		exitOnError("unable to get current work dir", err)
	}
	helpers.ProjectDir = dir
	logrus.WithField("ProjectDir", helpers.ProjectDir).Debug("using cwd")

	file, err := rconfig.LoadFile()
	if err != nil {
		exitOnError("file read", err)
	}

	config, err := rconfig.ParseConfig(file)
	if err != nil {
		exitOnError("file parse", err)
	}

	ctx := context.TODO() // @todo ^^
	dockerClient, err := helpers.NewDockerClient()
	if err != nil {
		exitOnError("docker connection", err)
	}

	if *cmdStop {
		containers, listContainersBeforeStopErr := helpers.ListContainers(ctx, dockerClient)
		if listContainersBeforeStopErr != nil {
			exitOnError("list problem", listContainersBeforeStopErr)
		}
		logic.StopContainers(ctx, dockerClient, containers)
		os.Exit(0)
	}

	if !*cmdList && !*watch {
		os.Exit(0)
	}

	for {
		containers, errList := helpers.ListContainers(ctx, dockerClient)
		if errList != nil {
			logrus.WithField("error", errList).Error("refreshing container list")
		}
		countContainers := logic.CountContainers(containers)
		logic.MngContainers(ctx, config, dockerClient, countContainers)

		if !*watch {
			break
		}

		time.Sleep(3 * time.Second) // @todo a config option would be nice
		logrus.Info("refreshing ^^")
	}
}

func exitOnError(message string, err error) {
	logrus.WithField("error", err).Fatal(message)
}
