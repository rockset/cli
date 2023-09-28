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

type SelectionString struct {
	Path       []PathElem
	ColumnName string
}

type Selector []SelectionString

func ParseSelectionString(s string) (Selector, error) {
	elems := strings.Split(s, ",")
	columns := make([]SelectionString, 0)

	for _, elem := range elems {
		segments := strings.Split(elem, ":")
		columnName := segments[0]
		pathString := segments[0]
		if len(segments) == 2 {
			pathString = segments[1]
		}

		path := strings.Split(pathString, ".")[1:]
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
					return make([]SelectionString, 0), err
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

		columns = append(columns, SelectionString{
			Path:       richPath,
			ColumnName: columnName,
		})
	}

	return columns, nil
}

func (r Selector) Headers() []string {
	headers := make([]string, len(r))
	for i, sel := range r {
		headers[i] = sel.ColumnName
	}
	return headers
}

func (r Selector) Fields(data any) ([]string, error) {
	fields := make([]string, len(r))
	for i, sel := range r {
		value, err := sel.Select(data)
		if err != nil {
			return nil, err
		}
		valueAsString, err := AnyAsString(value)
		if err != nil {
			return nil, err
		}
		fields[i] = valueAsString
	}
	return fields, nil
}

func findFieldByJsonTag(value reflect.Value, jsonTag string) (reflect.Value, error) {
	for i := 0; i < value.Type().NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get("json")
		jsonName := strings.Split(tag, ",")[0]
		if jsonName == jsonTag {
			return value.FieldByName(value.Type().Field(i).Name), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("could not find json tag %s in type %s", jsonTag, value.Type().Name())
}

func (r SelectionString) Select(obj any) (any, error) {
	cur := reflect.Indirect(reflect.ValueOf(obj))
	curPath := r.Path

	for len(curPath) > 0 {
		next, rest := curPath[0], curPath[1:]

		if !cur.IsValid() {
			return nil, fmt.Errorf("selector %s is not valid on type %s", r.ToString(), reflect.TypeOf(obj).Name())
		} else if cur.Kind() == reflect.Ptr && cur.IsNil() {
			return nil, nil
		} else if cur.IsZero() {
			return nil, nil
		}

		var err error
		cur, err = findFieldByJsonTag(reflect.Indirect(cur), next.FieldName)
		if err != nil {
			return nil, fmt.Errorf("selector %s is not valid on type %s: %s", r.ToString(), reflect.TypeOf(obj).Name(), err.Error())
		}

		if next.HasArraySelector {
			cur = cur.Index(next.ArrayIndex)
		}

		curPath = rest
	}

	if !cur.IsValid() {
		return nil, fmt.Errorf("selector %s is not valid on type %s", r.ToString(), reflect.TypeOf(obj).Name())
	}

	return cur.Interface(), nil
}

func (r SelectionString) ToString() string {
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
