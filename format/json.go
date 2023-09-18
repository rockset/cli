package format

import (
	"encoding/json"
	"io"
)

type JSON struct {
	out io.Writer
}

func NewJSONFormatter(out io.Writer, header bool) *JSON {
	return &JSON{out}
}

func (j JSON) Format(wide bool, i interface{}) error {
	return json.NewEncoder(j.out).Encode(i)
}

func (j JSON) FormatList(wide bool, items []interface{}) error {
	return json.NewEncoder(j.out).Encode(items)
}
