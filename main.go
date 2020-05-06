package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/rockset/cli/cmd"
	"os"
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
		red := color.New(color.FgRed).SprintFunc()
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", red("ERROR"), err)
		os.Exit(1)
	}
}
