package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "heroku",
	Aliases: []string{"h"},
	Short:   "generates a heroku.yml file for deploying a buffalo app to heroku with docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		// return deployContainer()
		f, err := os.Create("heroku.yml")
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = f.WriteString(yml)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	},
}

const yml = `build:
  docker:
    web: Dockerfile

release:
  command:
    - /bin/app migrate
  image: web
`
