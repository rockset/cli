package format

import (
	"github.com/olekukonko/tablewriter"
	"io"
)

type Table struct {
	Header bool
	out    *tablewriter.Table
}

func NewTableFormatter(out io.Writer, header bool) *Table {
	return &Table{
		Header: header,
		out:    tablewriter.NewWriter(out),
	}
}

func (t Table) Format(wide bool, item interface{}) error {
	f, err := StructFormatterFor(item)
	if err != nil {
		return err
	}

	if t.Header {
		t.out.SetHeader(f.Headers(wide))
	}
	t.out.Append(f.Fields(wide, item))
	t.out.Render()

	return nil
}

func (t Table) FormatList(wide bool, items []interface{}) error {
	if len(items) == 0 {
		t.out.Render()
		return nil
	}

	f, err := StructFormatterFor(items[0])
	if err != nil {
		return err
	}
	if t.Header {
		t.out.SetHeader(f.Headers(wide))
	}

	for _, i := range items {
		t.out.Append(f.Fields(wide, i))
	}
	t.out.Render()

	return nil
}
