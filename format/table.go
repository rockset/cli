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

func (t Table) Format(wide bool, selector Selector, i interface{}) error {
	if selector == nil {
		defaults, err := DefaultSelectorFor(i, wide)
		if err != nil {
			return err
		}
		selector = defaults
	}

	if t.Header {
		t.w.SetHeader([]string{"key", "value"})
	}

	for _, sel := range selector {
		value, err := sel.Select(i)
		if err != nil {
			return err
		}
		valueAsString, err := AnyAsString(value, sel.FieldFormatter)
		if err != nil {
			return err
		}
		t.w.Append([]string{sel.ColumnName, valueAsString})
	}

	t.w.Render()

	return nil
}

func (t Table) FormatList(wide bool, selector Selector, items []interface{}) error {
	if items == nil || len(items) == 0 {
		t.w.Render()
		return nil
	}

	if selector == nil {
		defaults, err := DefaultSelectorFor(items[0], wide)
		if err != nil {
			return err
		}
		selector = defaults
	}

	if t.Header {
		t.w.SetHeader(selector.Headers())
	}

	for _, item := range items {
		fields, err := selector.Fields(item)
		if err != nil {
			return err
		}
		t.w.Append(fields)
	}

	t.w.Render()

	return nil
}
