package filter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	// FieldNameSeparator is the separator for Rule fields.
	FieldNameSeparator = "."
)

// reflectPath represents a field path (e.g. address.country) as the indices of the fields (e.g. [2,1]) that can be used with reflect.Value.Field(i int).
type reflectPath []int

type fieldGetter struct {
	fieldTag string
	cache    fieldCache
}

// GetFieldValue returns the value of obj's field, specified by its path.
func (r *fieldGetter) GetFieldValue(path string, obj interface{}) (interface{}, error) {
	// root path case
	if path == "" {
		return obj, nil
	}

	if obj == nil {
		return nil, errors.New("cannot get a field of a nil object")
	}

	tElement := reflect.TypeOf(obj)

	rPath, ok := r.cache.Get(tElement, path)
	if !ok {
		var err error

		pathParts := strings.Split(path, FieldNameSeparator)

		rPath, err = r.getFieldPath(pathParts, tElement)
		if err != nil {
			return nil, err
		}

		r.cache.Set(tElement, path, rPath)
	}

	value := reflect.ValueOf(obj)
	for _, fieldIndex := range rPath {
		value = reflect.Indirect(value)
		value = value.Field(fieldIndex)
	}

	if !value.CanInterface() {
		return nil, fmt.Errorf("%s cannot be interfaced", value.Type())
	}

	return value.Interface(), nil
}

func (r *fieldGetter) getFieldPath(fieldNames []string, t reflect.Type) (reflectPath, error) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("fields of elements of type %s are not supported", t)
	}

	currentName := fieldNames[0]

	field, err := r.getStructField(t, currentName)
	if err != nil {
		return nil, err
	}

	fieldPath := field.Index

	if len(fieldNames) > 1 {
		subPath, err := r.getFieldPath(fieldNames[1:], field.Type)
		if err != nil {
			return nil, err
		}

		fieldPath = append(fieldPath, subPath...)
	}

	return fieldPath, nil
}

func (r *fieldGetter) getStructField(t reflect.Type, name string) (reflect.StructField, error) {
	if r.fieldTag == "" {
		field, ok := t.FieldByName(name)
		if !ok {
			return reflect.StructField{}, fmt.Errorf("struct %s does not have a field named %s", t, name)
		}

		return field, nil
	}

	field, ok := r.lookupFieldByTag(t, name)
	if !ok {
		return reflect.StructField{}, fmt.Errorf("struct %s does not have a field with %s tag value of %s", t, r.fieldTag, name)
	}

	return field, nil
}

func (r *fieldGetter) lookupFieldByTag(t reflect.Type, tagValue string) (reflect.StructField, bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		actualValue := field.Tag.Get(r.fieldTag)
		if actualValue == tagValue {
			return field, true
		}
	}

	return reflect.StructField{}, false
}
