package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/structs"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// availableCmd represents the available command
var setupCmd = &cobra.Command{
	Use:   "heroku",
	Short: "setup heroku for deploying with docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return setup.Run()
	},
}

var setup = Setup{}

func init() {
	setupCmd.Flags().StringVarP(&setup.AppName, "app-name", "a", "", "the name for the heroku app")
	setupCmd.Flags().StringVarP(&setup.Environment, "environment", "e", "production", "setting for the GO_ENV variable")
	setupCmd.Flags().StringVarP(&setup.Database, "database", "d", "hobby-dev", "level of postgres database to use. use empty string for no database")
	setupCmd.Flags().BoolVarP(&setup.SkipAuth, "skip-auth", "s", false, "skip authorization")
	RootCmd.AddCommand(setupCmd)
}

type Setup struct {
	AppName     string
	Environment string
	Database    string
	SkipAuth    bool
}

func (s Setup) Run() error {
	g := makr.New()
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			if _, err := exec.LookPath("heroku"); err != nil {
				if runtime.GOOS == "darwin" {
					if _, err := exec.LookPath("brew"); err == nil {
						c := exec.Command("brew", "install", "heroku")
						c.Stdin = os.Stdin
						c.Stderr = os.Stderr
						c.Stdout = os.Stdout
						return c.Run()
					}
				}
				return errors.New("heroku cli is not installed. https://devcenter.heroku.com/articles/heroku-cli")
			}
			fmt.Println("--> heroku cli is installed")
			return nil
		},
	})
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			c := exec.Command("heroku", "plugins")
			b, err := c.CombinedOutput()
			if err != nil {
				return errors.WithStack(err)
			}
			if !bytes.Contains(b, []byte("heroku-container-registry")) {
				c := exec.Command("heroku", "plugins:install", "heroku-container-registry")
				c.Stdin = os.Stdin
				c.Stderr = os.Stderr
				c.Stdout = os.Stdout
				return c.Run()
			}
			fmt.Println("--> heroku-container-registry plugin is installed")
			return nil
		},
	})
	if !s.SkipAuth {
		g.Add(makr.NewCommand(exec.Command("heroku", "login")))
		g.Add(makr.NewCommand(exec.Command("heroku", "container:login")))
	}
	g.Add(makr.NewCommand(exec.Command("heroku", "create", s.AppName)))
	g.Add(makr.NewCommand(exec.Command("heroku", "config:set", fmt.Sprintf("GO_ENV=%s", s.Environment))))
	if s.Database != "" {
		g.Add(makr.NewCommand(exec.Command("heroku", "addons:create", fmt.Sprintf("heroku-postgresql:%s", s.Database))))
	}
	g.Add(makr.NewCommand(exec.Command("heroku", "container:push", "web")))
	return g.Run(".", structs.Map(s))
}
