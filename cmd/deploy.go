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
	if _, err := os.Stat("heroku.yml"); err == nil {
		fmt.Println("found a heroku.yml file; deploying with that")
		return pushHerokuYml()
	}
	if err := pushContainer(); err != nil {
		return errors.WithStack(err)
	}

	if err := releaseContainer(); err != nil {
		return errors.WithStack(err)
	}

	return runMigrations()
}

func pushHerokuYml() error {
	for _, f := range []string{"Dockerfile", "heroku.yml"} {
		c := exec.Command("git", "add", f)
		fmt.Println(strings.Join(c.Args, " "))
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			fmt.Println(err)
		}
	}
	c := exec.Command("git", "commit", "-m", "files for buffalo-heroku")
	fmt.Println(strings.Join(c.Args, " "))
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		fmt.Println(err)
	}

	c = exec.Command("git", "push", "heroku", "master")
	fmt.Println(strings.Join(c.Args, " "))
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
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

	return nil
}
func init() {
	herokuCmd.AddCommand(deployCmd)
}
