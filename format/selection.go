package format

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type PathElem struct {
	FieldName string

	HasArraySelector bool
	ArrayIndex       int
}

type Selector []FieldSelection

type FieldSelection struct {
	Path           []PathElem
	ColumnName     string
	FieldFormatter FieldFormatter
}

func NewFieldSelection(columnName string, path ...string) FieldSelection {
	fs := FieldSelection{
		ColumnName: columnName,
	}
	for _, p := range path {
		fs.Path = append(fs.Path, PathElem{FieldName: p})
	}
	return fs
}

func (s Selector) String() string {
	var list []string

	for _, f := range s {
		var sb strings.Builder
		sb.WriteString(f.ColumnName)
		sb.WriteString(":")
		for _, p := range f.Path {
			sb.WriteString(".")
			sb.WriteString(p.FieldName)
		}
		if f.FieldFormatter != nil {
			sb.WriteString(":")
			sb.WriteString(f.FieldFormatter.Name())
		}
		list = append(list, sb.String())
	}

	return strings.Join(list, ",")
}

func ParseSelectionString(s string) (Selector, error) {
	elems := strings.Split(s, ",")
	columns := make([]FieldSelection, 0)

	for _, elem := range elems {
		segments := strings.Split(elem, ":")
		fs := FieldSelection{ColumnName: segments[0]}

		pathString := segments[0]
		if len(segments) == 2 {
			pathString = segments[1]
		} else if len(segments) == 3 {
			pathString = segments[1]
			ff, found := FieldFormatters[segments[2]]
			if !found {
				return nil, fmt.Errorf("field formatter %s not found", segments[2])
			}
			fs.FieldFormatter = ff
		}

		path := strings.Split(pathString, ".")[1:]
		if len(path) == 0 {
			continue
		}

		if pathString == "." {
			path = make([]string, 0)
		}

		richPath := make([]PathElem, 0)
		for _, v := range path {
			pathSplit := strings.Split(v, "[")
			field := pathSplit[0]
			if len(pathSplit) == 2 {
				res, err := strconv.ParseInt(strings.Trim(pathSplit[1], "[] "), 10, 64)
				if err != nil {
					return make([]FieldSelection, 0), err
				}
				richPath = append(richPath, PathElem{
					FieldName:        field,
					HasArraySelector: true,
					ArrayIndex:       int(res),
				})
			} else {
				richPath = append(richPath, PathElem{
					FieldName:        field,
					HasArraySelector: false,
					ArrayIndex:       0,
				})
			}
		}

		fs.Path = richPath
		columns = append(columns, fs)
	}

	if len(columns) == 0 {
		return columns, fmt.Errorf("selector '%s' invalid; double-check to ensure that your selectors begin with a period", s)
	}

	return columns, nil
}

func (s Selector) Headers() []string {
	headers := make([]string, len(s))
	for i, sel := range s {
		headers[i] = sel.ColumnName
	}
	return headers
}

func (s Selector) Fields(data any) ([]string, error) {
	fields := make([]string, len(s))
	for i, sel := range s {
		value, err := sel.Select(data)
		if err != nil {
			return nil, err
		}
		valueAsString, err := AnyAsString(value, sel.FieldFormatter)
		if err != nil {
			return nil, err
		}
		fields[i] = valueAsString
	}
	return fields, nil
}

func getJsonTagOfField(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	return strings.Split(tag, ",")[0]
}

func findFieldByJsonTag(value reflect.Value, jsonTag string) (reflect.Value, error) {
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("invalid selector, type %s is not a struct", value.Type())
	}
	for i := 0; i < value.Type().NumField(); i++ {
		tag := getJsonTagOfField(value.Type().Field(i))
		if tag == jsonTag {
			return value.FieldByName(value.Type().Field(i).Name), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("could not find json tag %s in type %s", jsonTag, value.Type().Name())
}

func (r FieldSelection) Select(obj any) (any, error) {
	cur := reflect.Indirect(reflect.ValueOf(obj))
	curPath := r.Path

	makeError := func(err error) error {
		typeName := reflect.TypeOf(obj).Name()
		possibleSelectors := GetPossibleSelectorsFor(obj)
		selector := r.ToString()

		prefix := fmt.Sprintf("selector %s is not valid on type %s", selector, typeName)
		suffix := fmt.Sprintf("possible selectors on type %s include: %s", typeName, strings.Join(possibleSelectors, "; "))
		if err == nil {
			return fmt.Errorf("%s - %s", prefix, suffix)
		} else {
			return fmt.Errorf("%s: %s - %s", prefix, err.Error(), suffix)
		}

	}

	for len(curPath) > 0 {
		next, rest := curPath[0], curPath[1:]

		if !cur.IsValid() {
			return nil, makeError(nil)
		} else if cur.Kind() == reflect.Ptr && cur.IsNil() {
			return nil, nil
		} else if cur.IsZero() {
			return nil, nil
		}

		var err error
		cur, err = findFieldByJsonTag(reflect.Indirect(cur), next.FieldName)
		if err != nil {
			return nil, makeError(err)
		}

		if next.HasArraySelector {
			if next.ArrayIndex < 0 || next.ArrayIndex >= cur.Len() {
				return nil, nil
			}
			cur = cur.Index(next.ArrayIndex)
		}

		curPath = rest
	}

	if !cur.IsValid() {
		return nil, makeError(nil)
	}

	return cur.Interface(), nil
}

func (r FieldSelection) ToString() string {
	path := make([]string, 0)
	for _, elem := range r.Path {
		path = append(path, elem.ToString())
	}
	return "." + strings.Join(path, ".")
}

func (r PathElem) ToString() string {
	if r.HasArraySelector {
		return fmt.Sprintf("%s[%d]", r.FieldName, r.ArrayIndex)
	}
	return r.FieldName
}

func GetPossibleSelectorsFor(data any) []string {
	type SelectorQueueObject struct {
		typ    reflect.Type
		prefix string
	}

	selectors := make([]string, 1)
	selectors[0] = "."
	queue := make([]SelectorQueueObject, 1)
	queue[0] = SelectorQueueObject{
		typ:    reflect.TypeOf(data),
		prefix: "",
	}

	for len(queue) > 0 {
		cur, rest := queue[0], queue[1:]
		queue = rest

		for i := 0; i < cur.typ.NumField(); i++ {
			field := cur.typ.Field(i)
			tag := getJsonTagOfField(field)
			if tag == "" {
				continue
			}
			myPrefix := cur.prefix + "." + tag
			selectors = append(selectors, myPrefix)

			if field.Type.Kind() == reflect.Struct {
				queue = append(queue, SelectorQueueObject{
					typ:    field.Type,
					prefix: myPrefix,
				})
			}
			if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				queue = append(queue, SelectorQueueObject{
					typ:    field.Type.Elem(),
					prefix: myPrefix,
				})
			}
		}

	}
	return selectors
}
