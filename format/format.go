package format

import (
	"fmt"
	"github.com/rockset/rockset-go-client/openapi"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
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

const (
	CSVFormat   Format = "csv"
	TableFormat        = "table"
)

func FormatterFor(out io.Writer, f Format, header bool) Formatter {
	switch f {
	case CSVFormat:
		return NewCSVFormat(out, header)
	case TableFormat:
		return NewTableFormatter(out, header)
	default:
		return nil
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
	case openapi.User:
		return UserFormatter, nil
	case openapi.Organization:
		return OrgFormatter, nil
	default:
		return StructFormatter{}, fmt.Errorf("no formatter for %T", t)
	}
}

type StructFormatter struct {
	headers []Header
}

func (s StructFormatter) filteredHeaders(wide bool) []Header {
	var hdrs []Header
	for _, h := range s.headers {
		if wide || !h.Wide {
			hdrs = append(hdrs, h)
		}
	}
	return hdrs
}

// Headers returns the list of header names
func (s StructFormatter) Headers(wide bool) []string {
	var headers []string
	for _, h := range s.headers {
		if wide || !h.Wide {
			headers = append(headers, h.DisplayName)
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
	// dereference pointer, if we got one
	if f.Kind() == reflect.Ptr {
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
		f = f.Elem()
	}

	for i := 0; i < f.Len(); i++ {
		item := f.Index(i)
		v := item.FieldByName("Id")
		if v.IsValid() {
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			a = append(a, v.String())
		} else {
			a = append(a, item.String())
		}

	}
	return strings.Join(a, ", ")
}

func headers(i interface{}) []string {
	switch t := i.(type) {
	case openapi.Collection:
		return []string{
			"Workspace",
			"Name",
			"Description",
			"RetentionSecs",
			"Status",
			"InsertOnly",
			"EnableExactlyOnceWrites",
			"CreatedBy",
			"CreatedAt",
		}
	case openapi.Organization:
		return []string{
			"DisplayName",
			"RocksetUser",
			"Id",
			"ExternalId",
			"CreatedAt",
			"DeletionScheduledAt",
			// TODO: needs a formatter
			// "Clusters",
		}
	case openapi.QueryLambda:
		return []string{
			"Workspace",
			"Name",
			"LastUpdatedBy",
			"LastUpdated",
			// TODO: needs a formatter
			//"LatestVersion",
			"VersionCount",
			// TODO: needs a formatter
			//"Collections",
		}
	case openapi.Workspace:
		return []string{
			"Name",
			"Description",
			"CollectionCount",
			"CreatedBy",
			"CreatedAt",
		}
	default:
		log.Fatalf("unknown type %T", t)
	}
	return nil
}