package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runMigrations() error {
	if _, err := os.Stat("./database.yml"); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		c := exec.CommandContext(ctx, "heroku", "run", "/bin/app", "migrate")
		fmt.Println(strings.Join(c.Args, " "))
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	}
	return nil
}
