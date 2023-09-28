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

func (c CSV) Format(wide bool, selector string, i interface{}) error {
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

	defer c.w.Flush()

	if c.Header {
		err := c.w.Write(selection.Headers())
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}
	fields, err := selection.Fields(i)
	if err != nil {
		return fmt.Errorf("failed to write csv: %w", err)
	}
	err = c.w.Write(fields)
	if err != nil {
		return fmt.Errorf("failed to write csv: %w", err)
	}
	return nil
}

func (c CSV) FormatList(wide bool, selector string, items []interface{}) error {
	if len(items) == 0 {
		return fmt.Errorf("no items in list")
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

	defer c.w.Flush()

	if c.Header {
		err := c.w.Write(selection.Headers())
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}

	for _, item := range items {
		fields, err := selection.Fields(item)
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
