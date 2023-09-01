package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"

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

	// TODO should this be done in a PersistentPreRun & PersistentPostRun instead?
	// fire off a go routine to get the latest version
	version := make(chan string, 1)
	versionCtx, tc := context.WithTimeout(ctx, time.Second)
	defer tc()
	go cmd.VersionCheck(versionCtx, version)

	root := cmd.NewRootCmd(Version)
	if err := root.ExecuteContext(ctx); err != nil {
		errorf := color.New(color.Bold, color.FgRed).FprintfFunc()
		errorf(os.Stderr, "\nERROR: %v\n", err)
		os.Exit(1)
	}

	// show a warning if there is a new version available
	if v := <-version; v != "" {
		fmt.Printf("\n%s\n", v)
	}
}
