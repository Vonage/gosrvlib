package filter

import (
	"fmt"
	"reflect"
	"regexp"
)

const (
	// TypeExact is a filter type that matches exactly the reference value.
	TypeExact = "exact"

	// TypeNotExact is a filter type that matches when the value is different from the reference value (opposite of TypeExact).
	TypeNotExact = "different"

	// TypeRegexp is a filter type that matches the value against a reference regular expression.
	// The reference value must be a regular expression that can compile
	// Works only with strings (anything else will evaluate to false).
	TypeRegexp = "regexp"
)

// Evaluator is the interface to provide functions for a filter type.
type Evaluator interface {
	// Evaluate determines if two given values match
	Evaluate(reference, actual interface{}) (bool, error)
}

func getRuleType(typeName string) (Evaluator, error) {
	switch typeName {
	case TypeExact:
		return &exact{}, nil
	case TypeNotExact:
		return &not{Opposite: &exact{}}, nil
	case TypeRegexp:
		return &evalRegexp{}, nil
	default:
		return nil, fmt.Errorf("type %s is not supported", typeName)
	}
}

type exact struct {
}

func (e *exact) Evaluate(reference, actual interface{}) (bool, error) {
	actual = convertValues(actual)
	reference = convertValues(reference)

	if actual == reference {
		return true, nil
	}

	if e.isNil(actual) && e.isNil(reference) {
		return true, nil
	}

	return false, nil
}

func (e *exact) isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	if (value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr) && value.IsNil() {
		return true
	}

	return false
}

// nolint: gocyclo
func convertValues(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return float64(value)
	case int8:
		return float64(value)
	case int16:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	case uint:
		return float64(value)
	case uint8:
		return float64(value)
	case uint16:
		return float64(value)
	case uint32:
		return float64(value)
	case uint64:
		return float64(value)
	case float32:
		return float64(value)
	default:
		return value
	}
}

type not struct {
	Opposite Evaluator
}

func (n *not) Evaluate(reference, actual interface{}) (bool, error) {
	res, err := n.Opposite.Evaluate(reference, actual)
	if err != nil {
		return false, fmt.Errorf("failed evaluating the opposite rule: %w", err)
	}

	return !res, nil
}

type evalRegexp struct {
	internal *regexp.Regexp
}

func (r *evalRegexp) Evaluate(reference, actual interface{}) (bool, error) {
	if r.internal == nil {
		err := r.compile(reference)
		if err != nil {
			return false, err
		}
	}

	actualStr, ok := actual.(string)
	if !ok {
		return false, nil
	}

	return r.internal.MatchString(actualStr), nil
}

func (r *evalRegexp) compile(ref interface{}) error {
	str, ok := ref.(string)
	if !ok {
		return fmt.Errorf("rule of type %s should have string values (got %v (%v))", TypeRegexp, ref, reflect.TypeOf(ref))
	}

	reg, err := regexp.Compile(str)
	if err != nil {
		return fmt.Errorf("failed compiling regexp: %w", err)
	}

	r.internal = reg

	return nil
}
