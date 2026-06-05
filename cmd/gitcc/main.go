// Package main is the entry point for the gitcc CLI application.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/IceflowRE/gitcc/v3/standalone/gitcc/cli"
)

func main() {
	c := cli.NewCli(nil)

	err := c.Execute()
	if err == nil {
		return
	}

	var exitErr *cli.ExitError

	if errors.As(err, &exitErr) {
		os.Exit(exitErr.Code)
	}
	fmt.Fprintln(os.Stderr, err)

	os.Exit(3) //nolint:revive
}
