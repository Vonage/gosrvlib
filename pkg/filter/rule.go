package filter

import (
	"fmt"
	"strings"
)

const (
	// TypeEqual is a filter type that matches exactly the reference value.
	TypeEqual = "equal"
	// TypeEqualAliasA is an alias for TypeEqual.
	TypeEqualAliasA = "="
	// TypeEqualAliasB is an alias for TypeEqual.
	TypeEqualAliasB = "=="

	// TypeNotEqual is a filter type that matches when the value is different from the reference value (opposite of TypeEqual).
	TypeNotEqual = "notequal"
	// TypeNotEqualAliasA is an alias for TypeNotEqual.
	TypeNotEqualAliasA = "!="
	// TypeNotEqualAliasB is an alias for TypeNotEqual.
	TypeNotEqualAliasB = "<>"

	// TypeRegexp is a filter type that matches the value against a reference regular expression.
	// The reference value must be a regular expression that can compile.
	// Works only with strings (anything else will evaluate to false).
	TypeRegexp = "regexp"

	// TypeLT is a filter type that matches when the value is less than reference.
	TypeLT = "lt"
	// TypeLTAliasA is an alias for TypeLT.
	TypeLTAliasA = "<"

	// TypeLTE is a filter type that matches when the value is less than or equal the reference.
	TypeLTE = "lte"
	// TypeLTEAliasA is an alias for TypeLTE.
	TypeLTEAliasA = "<="

	// TypeGT is a filter type that matches when the value is greater than reference.
	TypeGT = "gt"
	// TypeGTAliasA is an alias for TypeGT.
	TypeGTAliasA = ">"

	// TypeGTE is a filter type that matches when the value is greater than or equal the reference.
	TypeGTE = "gte"
	// TypeGTEAliasA is an alias for TypeGTE.
	TypeGTEAliasA = ">="
)

// Rule is an individual filter that can be evaluated against any value.
type Rule struct {
	// Field is a dot separated selector that is used to target a specific field of the evaluated value.
	//
	// * "Age" will select the Age field of a structure
	// * "Address.Country" will select the Country subfield of the Address structure
	// * "" will select the whole value (e.g. to filter a []string)
	Field string `json:"field"`

	// Type controls the evaluation to apply.
	// An invalid value will cause Evaluate() to return an error.
	// See the Type* constants of this package for valid values.
	Type string `json:"type"`

	// Value is the reference value to evaluate against.
	// Its type should be accepted by the chosen Type.
	Value interface{} `json:"value"`

	// eval is initialized at the first call to Evaluate() and stores the structure that evaluates the rule.
	eval Evaluator
}

// Evaluate returns whether the value matches the rule or not.
//
// Returns an error if the Type is invalid, a misconfiguration (e.g. invalid regexp) or the value is invalid (e.g. evaluating an int with a regexp).
func (r *Rule) Evaluate(value interface{}) (bool, error) {
	if r.eval == nil {
		var err error

		r.eval, err = r.getEvaluator()
		if err != nil {
			return false, err
		}
	}

	return r.eval.Evaluate(value), nil
}

func (r *Rule) getEvaluator() (Evaluator, error) {
	switch strings.ToLower(r.Type) {
	case TypeEqual, TypeEqualAliasA, TypeEqualAliasB:
		return newEqual(r.Value), nil
	case TypeNotEqual, TypeNotEqualAliasA, TypeNotEqualAliasB:
		return newNot(newEqual(r.Value)), nil
	case TypeRegexp:
		return newRegexp(r.Value)
	case TypeLT, TypeLTAliasA:
		return newLT(r.Value)
	case TypeLTE, TypeLTEAliasA:
		return newLTE(r.Value)
	case TypeGT, TypeGTAliasA:
		e, err := newLTE(r.Value)
		if err != nil {
			return nil, err
		}

		return newNot(e), nil
	case TypeGTE, TypeGTEAliasA:
		e, err := newLT(r.Value)
		if err != nil {
			return nil, err
		}

		return newNot(e), nil
	default:
		return nil, fmt.Errorf("type %s is not supported", r.Type)
	}
}
