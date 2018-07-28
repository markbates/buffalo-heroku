package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "heroku",
	Aliases: []string{"h"},
	Short:   "generates a heroku.yml file for deploying a buffalo app to heroku with docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return WriteHerokuYml()
	},
}

func WriteHerokuYml() error {
	f, err := os.Create("heroku.yml")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = f.WriteString(yml)
	if err != nil {
		return errors.WithStack(err)
	}
	b, err := ioutil.ReadFile("Dockerfile")
	if err != nil {
		return errors.WithStack(err)
	}
	if bytes.Contains(b, []byte("RUN apk add --no-cache curl")) {
		return nil
	}
	b = bytes.Replace(b, []byte("FROM alpine\n"), []byte("FROM alpine\nRUN apk add --no-cache curl\n"), 1)
	f, err = os.Create("Dockerfile")
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

const yml = `build:
  docker:
    web: Dockerfile

release:
  command:
    - /bin/app migrate
  image: web
`
