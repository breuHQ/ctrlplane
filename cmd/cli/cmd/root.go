// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "ctrlplane",
		Short: "ctrlplane is a multi stage release rollout engine with pre-emptive rollbacks.",
		Long: `
ctrlplane is a multi stage release rollout engine for cloud-native applications. It is designed to be used in 
conjunction with a CI/CD pipeline & near realtime application monitoring to provide a safe and reliable rollout 
process with rollbacks for microservices.

Currently, it only supports monorepo, but poly repo support is on the roadmap.

To learn more, visit https://breu.io/ctrlplane.
  `,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if err := shared.Service.InitCLI(); err != nil {
		println("ctrlplane not initialized, please do ctrlplane init")
	}
}