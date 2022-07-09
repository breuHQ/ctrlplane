package main

import (
	"log"
	"sync"

	tw "go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/temporal/activities"
	"go.breu.io/ctrlplane/internal/temporal/workflows"
)

var wait sync.WaitGroup

func init() {
	conf.Service.ReadConf()
	conf.Service.InitLogger()

	conf.EventStream.ReadConf()
	conf.Temporal.ReadConf()
	conf.Github.ReadConf()
	conf.DB.ReadConf()

	wait.Add(3)

	go func() {
		defer wait.Done()
		conf.DB.InitSession()
	}()

	go func() {
		defer wait.Done()
		conf.EventStream.InitConnection()
	}()

	go func() {
		defer wait.Done()
		conf.Temporal.InitClient()
	}()

	wait.Wait()

	conf.Logger.Info("Initializing Service ... Done")
}

func main() {
	defer conf.Temporal.Client.Close()

	queue := conf.Temporal.Queues.Webhooks
	options := tw.Options{}
	worker := tw.New(conf.Temporal.Client, queue, options)

	worker.RegisterWorkflow(workflows.OnGithubInstall)
	worker.RegisterActivity(activities.SaveGithubInstallation)

	err := worker.Run(tw.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}
