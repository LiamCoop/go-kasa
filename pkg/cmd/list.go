package cmd

import (
	"github.com/liamcoop/go-kasa/internal"
	"github.com/spf13/cobra"
)

// should run discover first thing when first command is run,
// load into some local storage.
// then list is just going to iterate over the devices in storage
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list kasa devices on your network",
	Long:  `List your available kasa smart devices from within your network`,
	RunE: func(cmd *cobra.Command, args []string) error {
		devices, err := internal.Discover(1, 1)
		if err != nil {
			return err
		}

		internal.FormatDiscover(devices)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
