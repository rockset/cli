package format

import (
	"strconv"
	"time"

	"github.com/rockset/rockset-go-client/openapi"
)

func collection(c openapi.Collection, wide bool) []string {
	r := "forever"
	if c.GetRetentionSecs() != 0 {
		r = (time.Second * time.Duration(c.GetRetentionSecs())).String()
	}

	fields := []string{
		c.GetName(),
		c.GetWorkspace(),
		c.GetDescription(),
		c.GetStatus(),
		c.GetCreatedBy(),
		c.GetCreatedAt(),
		r,
	}

	if wide {
		fields = append(fields, strconv.FormatInt(c.Stats.GetTotalSize(), 10))
	}

	return fields
}

func collectionHeader(wide bool) []string {
	headers := []string{
		"name",
		"workspace",
		"description",
		"status",
		"created by",
		"created at",
		"retention",
	}
	if wide {
		headers = append(headers, "size")
	}

	return headers
}
