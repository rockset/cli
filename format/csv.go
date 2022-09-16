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

func (c CSV) Format(wide bool, i interface{}) error {
	f, err := StructFormatterFor(i)
	if err != nil {
		return err
	}

	defer c.w.Flush()

	if c.Header {
		err := c.w.Write(f.Headers(wide))
		if err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}

	err = c.w.Write(f.Fields(wide, i))
	if err != nil {
		return fmt.Errorf("failed to write csv: %w", err)
	}
	return nil
}

func (c CSV) FormatList(wide bool, i []interface{}) error {
	return nil
}
