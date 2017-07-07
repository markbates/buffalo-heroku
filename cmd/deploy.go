package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy to heroku using docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("heroku", "container:push", "web")
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
