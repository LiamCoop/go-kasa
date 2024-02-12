package cmd

import (
	"strings"

	"github.com/liamcoop/go-kasa/internal"
	"github.com/spf13/cobra"
)

// need some way to check that it's a device for which this is an allowed action.

// usage: go-kasa set <IP> <on/off>
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set value for digital devices (on/off)",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := args[0]
		d, err := internal.DeviceLookup(ip)
		if err != nil {
			return err
		}
		command := strings.ToLower(args[1])
		var parsedCommand bool
		if command == "on" {
			parsedCommand = true
		}

		if command == "off" {
			parsedCommand = false
		}

		d.SetState(parsedCommand)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
