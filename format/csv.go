package format

import (
	"encoding/csv"
	"fmt"
	"io"
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

func (c CSV) Format(wide bool, selector Selector, i interface{}) error {
	if selector == nil {
		defaults, err := DefaultSelectorFor(i, wide)
		if err != nil {
			return err
		}
		selector = defaults
	}

	defer c.w.Flush()

	if c.Header {

		if err := c.w.Write([]string{"key", "value"}); err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}

	for _, sel := range selector {
		value, err := sel.Select(i)
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
		valueAsString, err := AnyAsString(value, sel.FieldFormatter)
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
		err = c.w.Write([]string{sel.ColumnName, valueAsString})
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}
	return nil
}

func (c CSV) FormatList(wide bool, selector Selector, items []interface{}) error {
	if len(items) == 0 {
		return fmt.Errorf("no items in list")
	}
	if selector == nil {
		defaults, err := DefaultSelectorFor(items[0], wide)
		if err != nil {
			return err
		}
		selector = defaults
	}

	defer c.w.Flush()

	if c.Header {
		err := c.w.Write(selector.Headers())
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}

	for _, item := range items {
		fields, err := selector.Fields(item)
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
		err = c.w.Write(fields)
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}

	return nil
}
