package filter

import (
	"fmt"
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

// Evaluate returns true if the value matches the rule or not.
//
// Returns an error if the Type is invalid, a misconfiguration (e.g. invalid regexp) or the value is invalid (e.g. evaluating an int with a regexp).
func (r *Rule) Evaluate(value interface{}) (bool, error) {
	if r.eval == nil {
		var err error

		r.eval, err = getRuleType(r.Type)
		if err != nil {
			return false, err
		}
	}

	match, err := r.eval.Evaluate(r.Value, value)
	if err != nil {
		return false, fmt.Errorf("failed evaluating the rule: %w", err)
	}

	return match, nil
}
