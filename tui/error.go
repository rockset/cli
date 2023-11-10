package tui

import (
	"errors"
	"fmt"
	"io"

	rockerr "github.com/rockset/rockset-go-client/errors"
)

func ShowError(out io.Writer, debug bool, err error) {
	var re rockerr.Error
	if errors.As(err, &re) {
		if re.GetType() == "" {
			_, _ = fmt.Fprintf(out, "%s\n", ErrorStyle.Render("Error:", err.Error()))
		} else {
			_, _ = fmt.Fprintf(out, "%s\n", RocksetStyle.Render("Rockset error:", err.Error()))
		}

		if id := re.GetTraceId(); id != "" && debug {
			_, _ = fmt.Fprintf(out, "%s\n", RocksetStyle.Render("Trace ID:", id))
		}
	} else {
		_, _ = fmt.Fprintf(out, "%s\n", ErrorStyle.Render("Error:", err.Error()))
	}

}
