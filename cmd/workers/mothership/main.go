// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package main

import (
	"os"
	"sync"

	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	wait sync.WaitGroup
)

func init() {
	shared.Service.ReadEnv()
	shared.Service.InitLogger()
	shared.EventStream.ReadEnv()
	shared.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	wait.Add(3)

	go func() {
		defer wait.Done()
		db.DB.InitSession()
	}()

	go func() {
		defer wait.Done()
		shared.EventStream.InitConnection()
	}()

	go func() {
		defer wait.Done()
		shared.Temporal.InitClient()
	}()

	wait.Wait()

	shared.Logger.Info("Initializing Service ... Done", "version", shared.Service.Version())
}

func main() {
	// graceful shutdown. see https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()
	defer func() { _ = shared.Logger.Sync() }()
	defer func() { _ = shared.EventStream.Drain() }()
	defer shared.Temporal.Client.Close()

	queue := shared.Temporal.Queues[shared.ProvidersQueue].GetName()
	options := worker.Options{}
	wrkr := worker.New(shared.Temporal.Client, queue, options)

	workflows := &github.Workflows{}

	wrkr.RegisterWorkflow(workflows.OnInstall)
	wrkr.RegisterWorkflow(workflows.OnPush)
	wrkr.RegisterActivity(&github.Activities{})

	err := wrkr.Run(worker.InterruptCh())

	if err != nil {
		exitcode = 1
		return
	}
}
