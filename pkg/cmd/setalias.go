package cmd

import (
	"github.com/liamcoop/go-kasa/internal"
	"github.com/spf13/cobra"
)

// usage go-kasa alias <alias>
var setAliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Alias ",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := args[0]
		alias := args[1]
		d, err := internal.DeviceLookup(ip)
		if err != nil {
			return err
		}

		err = d.SetAlias(alias)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(setAliasCmd)
}
