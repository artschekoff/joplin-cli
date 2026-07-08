package main

import (
	"fmt"
	"os"

	"github.com/artschekoff/joplin-cli/src/internal/cli"
)

var version = "dev"

func main() {
	cli.SetVersion(version)
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
