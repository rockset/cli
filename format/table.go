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

func (t Table) Format(wide bool, selector string, i interface{}) error {
	if selector == "" {
		defaults, err := DefaultSelectorFor(i, wide)
		if err != nil {
			return err
		}
		selector = defaults
	}
	selection, err := ParseSelectionString(selector)
	if err != nil {
		return err
	}

	if t.Header {
		t.w.SetHeader(selection.Headers())
	}

	fields, err := selection.Fields(i)
	if err != nil {
		return err
	}

	t.w.Append(fields)
	t.w.Render()

	return nil
}

func (t Table) FormatList(wide bool, selector string, items []interface{}) error {
	if items == nil || len(items) == 0 {
		t.w.Render()
		return nil
	}

	if selector == "" {
		defaults, err := DefaultSelectorFor(items[0], wide)
		if err != nil {
			return err
		}
		selector = defaults
	}

	selection, err := ParseSelectionString(selector)
	if err != nil {
		return err
	}

	if t.Header {
		t.w.SetHeader(selection.Headers())
	}

	for _, item := range items {
		fields, err := selection.Fields(item)
		if err != nil {
			return err
		}
		t.w.Append(fields)
	}

	t.w.Render()

	return nil
}
