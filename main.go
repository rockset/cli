package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	rockerr "github.com/rockset/rockset-go-client/errors"

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
				errorf(fmt.Sprintf("program crash: %v", r))
				// TODO log message about the panic being sent to sentry
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
		// TODO allow users to override the error reporting
		// TODO log a message that we sent the error
		if !errors.Is(err, context.Canceled) {
			// we don't capture errors from the Rockset API
			var re rockerr.Error
			if errors.As(err, &re) {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", tui.RocksetStyle.Render("Rockset error:", err.Error()))
			} else {
				// this captures usage errors too, as there is no way to distinguish it from other errors
				sentry.CaptureException(err)
				errorf(err.Error())
			}
		}

		os.Exit(1)
	}

	// show a warning if there is a new version available, but on stderr as it will show up in the
	// rockset completion output otherwise
	if v := <-version; v != "" {
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", v)
	}
}

func errorf(msg string) {
	// bold too?
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", tui.ErrorStyle.Render("ERROR:", msg))
}
