package format

import (
	"github.com/olekukonko/tablewriter"
	"io"
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

func (t Table) Format(wide bool, i interface{}) error {
	f, err := StructFormatterFor(i)
	if err != nil {
		return err
	}

	if t.Header {
		t.w.SetHeader(f.Headers(wide))
	}
	t.w.Append(f.Fields(wide, i))
	t.w.Render()

	return nil
}

func (t Table) FormatList(wide bool, items []interface{}) error {
	if items == nil || len(items) == 0 {
		t.w.Render()
	}

	f, err := StructFormatterFor(items[0])
	if err != nil {
		return err
	}
	if t.Header {
		t.w.SetHeader(f.Headers(wide))
	}

	for _, i := range items {
		t.w.Append(f.Fields(wide, i))
	}
	t.w.Render()

	return nil
}
