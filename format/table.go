package format

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rockset/cli/tui"
	"io"
)

type Table struct {
	Header bool
	out    io.Writer
	table  *table.Table
}

func NewTableFormatter(out io.Writer, header bool) *Table {
	return &Table{
		Header: header,
		out:    out,
		// TODO might need a --no-color flag to turn coloring off
		table: tui.NewTable(out),
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
		t.table.Headers("key", "value")
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
		t.table.Row(sel.ColumnName, valueAsString)
		//t.w.Append([]string{sel.ColumnName, valueAsString})
	}

	_, _ = fmt.Fprintln(t.out, t.table.Render())

	return nil
}

func (t Table) FormatList(wide bool, selector Selector, items []interface{}) error {
	if items == nil || len(items) == 0 {
		// TODO what is a better message?
		t.table.Headers("No rows")
		_, _ = fmt.Fprintln(t.out, t.table.Render())
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
		t.table.Headers(selector.Headers()...)
	}

	for _, item := range items {
		fields, err := selector.Fields(item)
		if err != nil {
			return err
		}

		t.table.Row(fields...)
	}

	_, _ = fmt.Fprintln(t.out, t.table.Render())

	return nil
}
