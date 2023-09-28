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

	if len(columns) == 0 {
		return columns, fmt.Errorf("selector '%s' invalid; double-check to ensure that your selectors begin with a period", s)
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

func (r SelectionString) Select(obj any) (any, error) {
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
