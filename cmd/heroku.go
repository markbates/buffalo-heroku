package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// herokuCmd represents the heroku command
var herokuCmd = &cobra.Command{
	Use:     "heroku",
	Aliases: []string{"h"},
	Short:   "Tools for deploying Buffalo to Heroku",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("heroku called")
	},
}

func init() {
	RootCmd.AddCommand(herokuCmd)
}
