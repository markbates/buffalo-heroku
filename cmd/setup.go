package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/structs"
	"github.com/gobuffalo/makr"
	"github.com/markbates/going/randx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// availableCmd represents the available command
var setupCmd = &cobra.Command{
	Use:     "setup",
	Aliases: []string{"s"},
	Short:   "setup heroku for deploying with docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if setup.Free {
			setup.DynoType = "free"
			setup.Database = "hobby-dev"
		}
		if setup.Interactive {
			err := Interactive()
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return setup.Run()
	},
}

func Interactive() error {
	qs := []*survey.Question{
		{
			Name: "AppName",
			Prompt: &survey.Input{
				Message: "What would you like to name this app on Heroku?",
				Help:    "A blank response will let Heroku generate a name.",
			},
		},
		{
			Name: "DynoType",
			Prompt: &survey.Select{
				Message: "Choose a dyno level.",
				Options: dynoLevels,
				Default: setup.DynoType,
			},
		},
		{
			Name: "Database",
			Prompt: &survey.Select{
				Message: "Choose a PostgreSQL database level.",
				Options: pgLevels,
				Default: setup.Database,
			},
		},
		{
			Name: "EmailProvider",
			Prompt: &survey.Select{
				Message: "Choose an Email provider level.",
				Options: []string{"none", "sendgrid:starter"},
				Default: "none",
			},
		},
		{
			Name: "RedisProvider",
			Prompt: &survey.Select{
				Message: "Choose a Redis provider level.",
				Options: []string{"none", "heroku-redis:hobby-dev"},
				Default: "none",
			},
		},
	}
	err := survey.Ask(qs, &setup)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var dynoLevels = []string{"free", "hobby", "standard-1x", "standard-2x"}
var pgLevels = []string{"hobby-dev", "hobby-basic", "standard-0"}
var setup = Setup{}

func init() {
	setupCmd.Flags().StringVarP(&setup.AppName, "app-name", "a", "", "the name for the heroku app")
	setupCmd.Flags().StringVarP(&setup.Environment, "environment", "e", "production", "setting for the GO_ENV variable")
	setupCmd.Flags().StringVarP(&setup.Database, "database", "d", "hobby-basic", fmt.Sprintf("level of postgres database to use. use empty string for no database %s", pgLevels))
	setupCmd.Flags().StringVar(&setup.EmailProvider, "email", "sendgrid:starter", "email provider to use. use empty string for no database")
	setupCmd.Flags().StringVar(&setup.RedisProvider, "redis", "heroku-redis:hobby-dev", "redis provider to use. use empty string for no database")
	setupCmd.Flags().StringVarP(&setup.DynoType, "dyno-type", "t", "hobby", fmt.Sprintf("type of heroku dynos %s", dynoLevels))
	setupCmd.Flags().BoolVar(&setup.Auth, "auth", false, "perform authorization")
	setupCmd.Flags().BoolVarP(&setup.Free, "free", "f", false, "use only free resources")
	setupCmd.Flags().BoolVarP(&setup.Interactive, "interactive", "i", false, "use the interactive mode")
	herokuCmd.AddCommand(setupCmd)
}

//Setup holds the settings to setup the app
type Setup struct {
	AppName       string
	Environment   string
	Database      string
	Auth          bool
	DynoType      string
	EmailProvider string
	RedisProvider string
	Free          bool
	Interactive   bool
}

//Run runs setup steps
func (s Setup) Run() error {
	g := makr.New()
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return validateGit()
		},
	})
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return installHerokuCLI()
		},
	})

	if s.Auth {
		g.Add(makr.NewCommand(exec.Command("heroku", "login")))
		g.Add(makr.NewCommand(exec.Command("heroku", "container:login")))
	}
	g.Add(makr.NewCommand(exec.Command("heroku", "create", s.AppName)))
	g.Add(makr.NewCommand(exec.Command("heroku", "config:set", fmt.Sprintf("GO_ENV=%s", s.Environment))))
	g.Add(makr.NewCommand(exec.Command("heroku", "config:set", fmt.Sprintf("SESSION_SECRET=%s", randx.String(100)))))
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return initializeHostVar()
		},
	})
	if s.Database != "" {
		g.Add(makr.NewCommand(exec.Command("heroku", "addons:create", fmt.Sprintf("heroku-postgresql:%s", s.Database))))
	}
	if s.EmailProvider != "" && s.EmailProvider != "none" {
		g.Add(makr.NewCommand(exec.Command("heroku", "addons:create", s.EmailProvider)))
		if strings.Contains(s.EmailProvider, "sendgrid") {
			g.Add(makr.Func{
				Runner: func(root string, data makr.Data) error {
					return setupSendgrid()
				},
			})
		}
	}
	if s.RedisProvider != "" && s.RedisProvider != "none" {
		g.Add(makr.NewCommand(exec.Command("heroku", "addons:create", s.RedisProvider)))
	}
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return deployContainer()
		},
	})

	g.Add(makr.NewCommand(exec.Command("heroku", "dyno:type", s.DynoType)))
	g.Add(makr.NewCommand(exec.Command("heroku", "config")))
	g.Add(makr.NewCommand(exec.Command("heroku", "open")))
	return g.Run(".", structs.Map(s))
}

func setupSendgrid() error {
	cmd := exec.Command("heroku", "config:set", "SMTP_HOST=smtp.sendgrid.net", "SMTP_PORT=465")
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}

	cmd = exec.Command("heroku", "config:get", "SENDGRID_USERNAME")
	b, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}

	cmd = exec.Command("heroku", "config:set", fmt.Sprintf("SMTP_USER=%s", strings.TrimSpace(string(b))))
	b, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}

	cmd = exec.Command("heroku", "config:get", "SENDGRID_PASSWORD")
	b, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}

	cmd = exec.Command("heroku", "config:set", fmt.Sprintf("SMTP_PASSWORD=%s", strings.TrimSpace(string(b))))
	b, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}

	cmd = exec.Command("heroku", "config:set", "SMTP_PORT=25")
	b, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}

	return nil
}

func installHerokuCLI() error {
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
}
