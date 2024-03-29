package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/rockset/cli/flag"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/tui"
)

const publicDsn = "___PUBLIC_DSN___"

var dsn = publicDsn

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	sigs := make(chan os.Signal, 1)

	go func() {
		sig := <-sigs
		fmt.Printf("received signal %s, cancelling\n", sig)
		cancel()
		sig = <-sigs
		fmt.Printf("received second signal %s, exiting\n", sig)
		os.Exit(1)
	}()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Release:          Version,
		TracesSampleRate: 1.0,
	}); err != nil {
		// TODO log that we failed to init sentry
	}

	// TODO should this be done in a PersistentPreRun & PersistentPostRun instead?
	// fire off a go routine to get the latest version
	version := make(chan string, 1)
	versionCtx, tc := context.WithTimeout(ctx, time.Second)

	defer func() {
		if r := recover(); r != nil {
			if dsn == publicDsn {
				_, _ = fmt.Fprintf(os.Stderr, "panic: %v\n", r)
				_, _ = fmt.Fprintf(os.Stderr, "%s", string(debug.Stack()))
			} else {
				sentry.CurrentHub().Recover(r)
				// TODO log message about the panic being sent to sentry
				_, _ = fmt.Fprintf(os.Stderr, "%s %v\n", tui.ErrorStyle.Render("program crash:"), r)
			}
			os.Exit(1)
		}
		tc()

		sentry.Flush(2 * time.Second)
	}()

	// kick off a version check in the background that will show up at the end of the run
	go cmd.GithubVersionCheck(versionCtx, version)

	root := cmd.NewRootCmd(Version)
	if err := root.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			// TODO allow users to override the error reporting
			// TODO log a message that we sent the error
			// TODO this captures usage errors too, as there is no way to distinguish them from other errors
			sentry.CaptureException(err)

			dbg, _ := root.PersistentFlags().GetBool(flag.Debug)
			tui.ShowError(os.Stderr, dbg, err)
		}

		os.Exit(1)
	}

	// show a warning if there is a new version available, but on stderr as it will show up in the
	// rockset completion output otherwise
	if v := <-version; v != "" {
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", v)
	}
}
