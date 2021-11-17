package format

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/rockset/rockset-go-client/openapi"
)

type CSV struct {
	Header bool
	w      *csv.Writer
}

func NewCSVFormat(out io.Writer, header bool) CSV {
	return CSV{
		Header: header,
		w:      csv.NewWriter(out),
	}
}

func (c CSV) Workspace(ws openapi.Workspace) {
	if c.Header {
		err := c.w.Write(workspaceHeader())
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}

	err := c.w.Write(workspace(ws))
	if err != nil {
		log.Printf("failed to write csv: %v", err)
	}
	c.w.Flush()
}

func (c CSV) Workspaces(list []openapi.Workspace) {
	if c.Header {
		err := c.w.Write(workspaceHeader())
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}

	for _, ws := range list {
		err := c.w.Write(workspace(ws))
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}
	c.w.Flush()
}

func (c CSV) Collection(coll openapi.Collection, wide bool) {
	if c.Header {
		err := c.w.Write(collectionHeader(wide))
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}

	err := c.w.Write(collection(coll, wide))
	if err != nil {
		log.Printf("failed to write csv: %v", err)
	}
	c.w.Flush()
}

func (c CSV) Collections(list []openapi.Collection, wide bool) {
	if c.Header {
		err := c.w.Write(collectionHeader(wide))
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}

	for _, l := range list {
		err := c.w.Write(collection(l, wide))
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}
	c.w.Flush()
}

func (c CSV) Users(list []openapi.User) {
	if c.Header {
		err := c.w.Write(userHeader())
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}

	for _, u := range list {
		err := c.w.Write(user(u))
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}
	c.w.Flush()
}

func (c CSV) User(u openapi.User) {
	if c.Header {
		err := c.w.Write(userHeader())
		if err != nil {
			log.Printf("failed to write csv: %v", err)
		}
	}

	err := c.w.Write(user(u))
	if err != nil {
		log.Printf("failed to write csv: %v", err)
	}
	c.w.Flush()
}
