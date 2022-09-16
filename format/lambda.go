package format

import (
	"strconv"

	"github.com/rockset/rockset-go-client/openapi"
)

func lambda(ws openapi.QueryLambda) []string {
	return []string{
		ws.GetWorkspace(),
		ws.GetName(),
		ws.GetLastUpdatedBy(),
		ws.GetLastUpdated(),
		strconv.FormatInt(int64(ws.GetVersionCount()), 10),
	}
}

func lambdaHeader() []string {
	return []string{
		"workspace",
		"name",
		"updated by",
		"updated at",
		"version count",
	}
}
