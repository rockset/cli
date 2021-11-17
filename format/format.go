package format

import (
	"fmt"
	"io"

	"github.com/rockset/rockset-go-client/openapi"
)

type Formatter interface {
	Collection(c openapi.Collection, wide bool)
	Collections(list []openapi.Collection, wide bool)
	Workspace(ws openapi.Workspace)
	Workspaces(list []openapi.Workspace)
	Users(list []openapi.User)
	User(user openapi.User)
}

func FormatterFor(out io.Writer, f string, header bool) (Formatter, error) {
	switch f {
	case "csv":
		return NewCSVFormat(out, header), nil
	case "table":
		return NewTableFormatter(out, header), nil
	default:
		return nil, fmt.Errorf("unknown formatter: %s", f)
	}
}
