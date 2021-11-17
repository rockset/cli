package format

import (
	"io"
	"log"

	"github.com/olekukonko/tablewriter"

	"github.com/rockset/rockset-go-client/openapi"
)

type Table struct {
	Header bool
	w      *tablewriter.Table
}

func NewTableFormatter(out io.Writer, header bool) *Table {
	return &Table{
		Header: header,
		w:      tablewriter.NewWriter(out),
	}
}

func (t Table) Workspace(ws openapi.Workspace) {
	t.foo(ws, true)
}

func (t Table) Workspaces(list []openapi.Workspace) {
	if t.Header {
		t.w.SetHeader(workspaceHeader())
	}
	for _, ws := range list {
		t.w.Append(workspace(ws))
	}
	t.w.Render()
}

func (t Table) Collection(c openapi.Collection, wide bool) {
	t.foo(c, wide)
}

func (t Table) Collections(list []openapi.Collection, wide bool) {
	if t.Header {
		t.w.SetHeader(collectionHeader(wide))
	}
	for _, c := range list {
		t.w.Append(collection(c, wide))
	}
	t.w.Render()
}

func (t Table) User(user openapi.User) {
	t.foo(user, true)
}

func (t Table) Users(list []openapi.User) {
	if t.Header {
		t.w.SetHeader(userHeader())
	}
	for _, u := range list {
		t.w.Append(user(u))
	}
	t.w.Render()
}

func (t Table) foo(i interface{}, wide bool) {
	var header []string
	var row []string

	switch u := i.(type) {
	case openapi.User:
		header = userHeader()
		row = user(u)
	case openapi.Collection:
		header = collectionHeader(wide)
		row = collection(u, wide)
	default:
		log.Panicf("unknown type %T: %+v", i, i)
	}

	if t.Header {
		t.w.SetHeader(header)
	}
	t.w.Append(row)
	t.w.Render()
}
