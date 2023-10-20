package format

import (
	"encoding/json"
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
	Format(wide bool, selector Selector, item interface{}) error
	FormatList(wide bool, selector Selector, items []interface{}) error
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

func DefaultSelectorFor(f any, wide bool) (Selector, error) {
	sel, err := defaultSelectorFor(f)
	if err != nil {
		return nil, err
	}
	if sel.Wide != nil && wide {
		return sel.Wide, nil
	}
	return sel.Normal, nil
}

func defaultSelectorFor(f any) (DefaultSelector, error) {
	switch t := f.(type) {
	case openapi.Alias:
		return AliasDefaultSelector, nil
	case openapi.ApiKey:
		return ApiKeyDefaultSelector, nil
	case openapi.User:
		return UserDefaultSelector, nil
	case openapi.Organization:
		return OrganizationDefaultSelector, nil
	case openapi.Workspace:
		return WorkspaceDefaultSelector, nil
	case openapi.Collection:
		return CollectionDefaultSelector, nil
	case openapi.Integration:
		return IntegrationDefaultSelector, nil
	case openapi.QueryInfo:
		return QueryDefaultSelector, nil
	case openapi.QueryLambda:
		return QueryLambdaDefaultSelector, nil
	case openapi.QueryLambdaVersion:
		return QueryLambdaVersionDefaultSelector, nil
	case openapi.QueryLambdaTag:
		return QueryLambdaTagDefaultSelector, nil
	case openapi.View:
		return ViewDefaultSelector, nil
	case openapi.VirtualInstance:
		return VirtualInstanceDefaultSelector, nil
	case openapi.Role:
		return RoleDefaultSelector, nil
	default:
		return DefaultSelector{}, fmt.Errorf("no formatter for %T", t)
	}
}

type DefaultSelector struct {
	Normal Selector
	Wide   Selector
}

func valueAsString(value reflect.Value, ff FieldFormatter) (string, error) {
	k := value.Kind()
	switch k {
	case reflect.String:
		return value.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	case reflect.Int64, reflect.Int32:
		if ff == nil {
			return strconv.FormatInt(value.Int(), 10), nil
		}
		return ff.FormatField(value.Int())
	case reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Slice, reflect.Array:
		a := make([]string, value.Len())
		for i := 0; i < value.Len(); i++ {
			x := value.Index(i)
			parsed, err := valueAsString(x, ff)
			if err != nil {
				a[i] = fmt.Sprintf("%v", x)
			} else {
				a[i] = parsed
			}
		}

		return strings.Join(a, ", "), nil
	case reflect.Ptr:
		if value.IsNil() {
			return "", nil
		}
		return valueAsString(reflect.Indirect(value), ff)
	case reflect.Struct:
		out, err := json.Marshal(value.Interface())
		if err != nil {
			return "", err
		}
		return string(out), nil
	case reflect.Invalid:
		return "NULL", nil
	default:
		return "", fmt.Errorf("unhandled kind %s", k)
	}
}

func AnyAsString(a any, ff FieldFormatter) (string, error) {
	return valueAsString(reflect.ValueOf(a), ff)
}
