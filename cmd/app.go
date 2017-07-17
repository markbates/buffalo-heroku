package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type AppInfo struct {
	App App `json:"app"`
}

type App struct {
	WebURL string `json:"web_url"`
}

func validateGit() error {
	c := exec.Command("git", "status")
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		return errors.Wrap(err, "must be a valid git application")
	}
	return nil
}

func initializeHostVar() error {
	cmd := exec.Command("heroku", "apps:info", "-j")
	bb := &bytes.Buffer{}
	cmd.Stdout = bb
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	ai := AppInfo{}
	err = json.NewDecoder(bb).Decode(&ai)
	if err != nil {
		return errors.WithStack(err)
	}
	cmd = exec.Command("heroku", "config:set", fmt.Sprintf("HOST=%s", ai.App.WebURL))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
