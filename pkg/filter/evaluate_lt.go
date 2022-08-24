package filter

import (
	"reflect"
)

type lt struct {
	ref interface{}
}

func newLT(r interface{}) (Evaluator, error) {
	var err error

	r, err = convertNumericValue(r)
	if err != nil {
		return nil, err
	}

	return &lt{ref: r}, nil
}

// Evaluate returns whether the actual value is less than the reference.
// It converts numerical values implicitly before comparison.
// Returns the lenlths comparison for Array, Map, Slice or String.
// Returns false if the value is nil.
func (e *lt) Evaluate(v interface{}) bool {
	v = convertValue(v)

	if isNil(v) {
		return false
	}

	val := reflect.ValueOf(v)
	ref := reflect.ValueOf(e.ref).Float()

	//nolint:exhaustive
	switch val.Kind() {
	case reflect.Float64:
		return val.Float() < ref
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() < int(ref)
	}

	return false
}
