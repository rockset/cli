package format

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"time"
)

var FieldFormatters = map[string]FieldFormatter{
	SizeName:  &SizeFormatter{},
	TimeSince: &TimeSinceFormatter{},
}

type FieldFormatter interface {
	FormatField(any) (string, error)
	Name() string
}

type SizeFormatter struct{}

const SizeName = "size"

func (f SizeFormatter) Name() string { return SizeName }
func (f SizeFormatter) FormatField(a any) (string, error) {
	switch a.(type) {
	case int64:
		return humanize.Bytes(uint64(a.(int64))), nil
	case uint64:
		return humanize.Bytes(a.(uint64)), nil
	default:
		return "", fmt.Errorf("%v can't be turned into a size", a)
	}
}

const TimeSince = "time_since"

type TimeSinceFormatter struct{}

func (f TimeSinceFormatter) Name() string { return TimeSince }
func (f TimeSinceFormatter) FormatField(a any) (string, error) {
	switch a.(type) {
	case int64:
		i := a.(int64)
		if i == 0 {
			return "never", nil
		}
		t := time.UnixMilli(i)
		return humanize.Time(t), nil
	default:
		return "", fmt.Errorf("%v (%T) can't be turned into time", a, a)
	}
}
