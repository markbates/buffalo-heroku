package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"d"},
	Short:   "deploy to heroku using docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return deployContainer()
	},
}

func deployContainer() error {
	if err := pushContainer(); err != nil {
		return errors.WithStack(err)
	}
	if err := releaseContainer(); err != nil {
		return errors.WithStack(err)
	}
	return runMigrations()
}

func pushContainer() error {
	c := exec.Command("heroku", "container:push", "web")
	fmt.Println(strings.Join(c.Args, " "))
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func releaseContainer() error {
	c := exec.Command("heroku", "container:release", "web")
	fmt.Println(strings.Join(c.Args, " "))
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return runMigrations()
}
func init() {
	herokuCmd.AddCommand(deployCmd)
}
