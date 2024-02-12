package cmd

import (
	"fmt"

	"github.com/liamcoop/go-kasa/internal"
	"github.com/spf13/cobra"
)

// usage go-kasa rules <IP>
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List rules",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := args[0]
		d, err := internal.DeviceLookup(ip)
		if err != nil {
			return err
		}

		rules, err := d.GetRules()
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", rules)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(rulesCmd)
}
