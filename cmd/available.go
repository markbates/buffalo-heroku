package cmd

import (
	"encoding/json"
	"os"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/spf13/cobra"
)

// availableCmd represents the available command
var availableCmd = &cobra.Command{
	Use:   "available",
	Short: "a list of available buffalo plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		plugs := plugins.Commands{
			{Name: setupCmd.Use, BuffaloCommand: "setup", Description: setupCmd.Short},
			{Name: deployCmd.Use, BuffaloCommand: "root", Description: deployCmd.Short},
		}
		return json.NewEncoder(os.Stdout).Encode(plugs)
	},
}

func init() {
	RootCmd.AddCommand(availableCmd)
}
