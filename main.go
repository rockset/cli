package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rockset/cli/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)

	go func() {
		sig := <-sigs
		fmt.Printf("received signal %s, cancelling\n", sig)
		cancel()
		sig = <-sigs
		fmt.Printf("received second signal %s, exiting\n", sig)
		os.Exit(1)
	}()

	root := cmd.NewRootCmd()
	if err := root.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\nERROR: %v\n", err)
		os.Exit(1)
	}
}
