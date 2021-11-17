package format

import (
	"strconv"

	"github.com/rockset/rockset-go-client/openapi"
)

func workspace(ws openapi.Workspace) []string {
	return []string{
		ws.GetName(),
		ws.GetDescription(),
		ws.GetCreatedBy(),
		ws.GetCreatedAt(),
		strconv.FormatInt(ws.GetCollectionCount(), 10),
	}
}

func workspaceHeader() []string {
	return []string{
		"name",
		"description",
		"created by",
		"created at",
		"collection count",
	}
}
