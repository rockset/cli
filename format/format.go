package format

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/rockset/rockset-go-client/openapi"
)

func ToInterfaceArray[T any](list []T) []interface{} {
	a := make([]interface{}, len(list))
	for i, l := range list {
		a[i] = l
	}
	return a
}

type Formatter interface {
	Format(wide bool, item interface{}) error
	FormatList(wide bool, items []interface{}) error
}

type Format string
type Formats []Format

func (f Formats) ToStringArray() []string {
	var strs = make([]string, len(f))
	for i, s := range f {
		strs[i] = string(s)
	}

	return strs
}

const (
	CSVFormat   Format = "csv"
	TableFormat Format = "table"
	JSONFormat  Format = "json"
)

var SupportedFormats = Formats{CSVFormat, JSONFormat, TableFormat}

func FormatterFor(out io.Writer, f Format, header bool) (Formatter, error) {
	switch f {
	case CSVFormat:
		return NewCSVFormat(out, header), nil
	case TableFormat:
		return NewTableFormatter(out, header), nil
	case JSONFormat:
		return NewJSONFormatter(out, header), nil
	default:
		return nil, fmt.Errorf("unknown formatter '%s', possible values are %s", f,
			strings.Join(SupportedFormats.ToStringArray(), ", "))
	}
}

type Header struct {
	FieldName   string
	DisplayName string
	Wide        bool
	FieldFn     func(string, interface{}) string
}

func (h Header) Field(i interface{}) string {
	if h.FieldFn == nil {
		return getFieldByName(h.FieldName, i)
	}

	return h.FieldFn(h.FieldName, i)
}

func StructFormatterFor[T any](f T) (StructFormatter, error) {
	var i interface{} = f
	switch t := i.(type) {
	case openapi.Alias:
		return AliasFormatter, nil
	case openapi.ApiKey:
		return APIKeyFormatter, nil
	case openapi.User:
		return UserFormatter, nil
	case openapi.Organization:
		return OrgFormatter, nil
	case openapi.Workspace:
		return WorkspaceFormatter, nil
	case openapi.Collection:
		return CollectionFormatter, nil
	case openapi.QueryInfo:
		return QueryInfoFormatter, nil
	case openapi.QueryLambda:
		return QueryLambdaFormatter, nil
	case openapi.QueryLambdaTag:
		return QueryLambdaTagFormatter, nil
	case openapi.View:
		return ViewFormatter, nil
	case openapi.VirtualInstance:
		return VirtualInstanceFormatter, nil
	default:
		return StructFormatter{}, fmt.Errorf("no formatter for %T", t)
	}
}

type StructFormatter struct {
	headers []Header
}

func (s StructFormatter) filteredHeaders(wide bool) []Header {
	var headers []Header
	for _, h := range s.headers {
		if wide || !h.Wide {
			headers = append(headers, h)
		}
	}
	return headers
}

// Headers returns the list of header names
func (s StructFormatter) Headers(wide bool) []string {
	var headers []string
	for _, h := range s.headers {
		if wide || !h.Wide {
			if h.DisplayName == "" {
				headers = append(headers, h.FieldName)
			} else {
				headers = append(headers, h.DisplayName)
			}
		}
	}
	return headers
}

func (s StructFormatter) Fields(wide bool, i interface{}) []string {
	var fields []string
	for _, h := range s.filteredHeaders(wide) {
		fields = append(fields, h.Field(i))
	}
	return fields
}

func getFieldByName(name string, i interface{}) string {
	v := reflect.Indirect(reflect.ValueOf(i))
	f := v.FieldByName(name)

	if !f.IsValid() {
		return fmt.Sprintf("unknown field name %s", name)
	}

	// dereference pointer, if we got one
	if f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return ""
		}
		f = f.Elem()
	}

	switch k := f.Kind(); k {
	case reflect.String:
		return f.String()
	case reflect.Bool:
		return strconv.FormatBool(f.Bool())
	case reflect.Int64, reflect.Int32:
		return strconv.FormatInt(f.Int(), 10)
	case reflect.Uint64:
		return strconv.FormatUint(f.Uint(), 10)
	case reflect.Slice:
		a := make([]string, f.Len())
		for i := 0; i < f.Len(); i++ {
			x := f.Index(i)
			a[i] = fmt.Sprintf("%v", x)
		}

		return strings.Join(a, ", ")
	default:
		return fmt.Sprintf("[%T: unhandled kind %s for field %s]", i, k, name)
	}
}

func getArrayFieldByName(name string, i interface{}) string {
	var a []string
	v := reflect.Indirect(reflect.ValueOf(i))
	f := v.FieldByName(name)

	// dereference pointer, if we got one
	if f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return ""
		}
		f = f.Elem()
	}
	if !f.IsValid() {
		return "not valid"
	}

	for i := 0; i < f.Len(); i++ {
		item := f.Index(i)
		a = append(a, item.String())
	}
	return strings.Join(a, ", ")
}

func getArrayStructFieldByName(structFieldName string) func(string, any) string {
	return func(fieldName string, item any) string {

		var items []string
		v := reflect.Indirect(reflect.ValueOf(item))
		f := v.FieldByName(fieldName)

		// dereference pointer, if we got one
		if f.Kind() == reflect.Ptr {
			if f.IsNil() {
				return ""
			}
			f = f.Elem()
		}
		if !f.IsValid() {
			return "not valid"
		}

		for i := 0; i < f.Len(); i++ {
			str := f.Index(i)
			if str.Kind() == reflect.Ptr {
				str = f.Elem()
			}
			items = append(items, getFieldByName(structFieldName, str.Interface()))
		}
		return strings.Join(items, ", ")
	}
}

func getStructFieldByName(name string) func(string, any) string {
	return func(s string, a any) string {
		v := reflect.Indirect(reflect.ValueOf(a))
		f := v.FieldByName(s)

		if f.Kind() == reflect.Ptr {
			if f.IsNil() {
				return ""
			}
			f = f.Elem()
		}
		if !f.IsValid() {
			return "not valid"
		}

		sf := f.FieldByName(name)
		if sf.Kind() == reflect.Ptr {
			if sf.IsNil() {
				return ""
			}
			sf = sf.Elem()
		}

		return sf.String()
	}
}
