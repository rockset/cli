package main

import (
	"fmt"
	"github.com/rockset/cli/cmd"
	"os"
)

func main() {
	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
