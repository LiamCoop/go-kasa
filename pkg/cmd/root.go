package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "go-kasa <command> <subcommand> [flags]",
	Short: "Go-Kasa CLI",
	Long:  `Manage your kasa smart devices from within your network via CLI.`,
}

func init() {
	cobra.OnInitialize()
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
