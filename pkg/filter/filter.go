// Package filter provides generic filtering capabilities for struct slices.
package filter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

// Processor is the interface to provide subtractive functions.
type Processor interface {
	// Apply filters the slice to remove elements not matching the defined rules.
	// The slice parameter must be a pointer to a slice and is filtered *in place*.
	Apply(rules [][]Rule, slice interface{}) error
}

// GetFilter returns the "filter" query parameter from a *url.URL.
func GetFilter(u *url.URL) string {
	return u.Query().Get("filter")
}

// ParseRules parses and returns a [][]Rule from its JSON representation.
func ParseRules(s string) ([][]Rule, error) {
	var r [][]Rule
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return nil, fmt.Errorf("failed unmarshaling rules: %w", err)
	}

	return r, nil
}

// New returns a new Processor with the rules and the given options.
//
// The first level of rules is matched with an AND operator and the second level with an OR.
//
// "[a,[b,c],d]" evaluates to "a AND (b OR c) AND d".
func New(opts ...Option) (Processor, error) {
	p := processor{}

	for _, opt := range opts {
		if err := opt(&p); err != nil {
			return nil, err
		}
	}

	return &p, nil
}

type processor struct {
	fields fieldGetter
}

func (p *processor) Apply(rules [][]Rule, slicePtr interface{}) error {
	if len(rules) == 0 {
		return nil
	}

	vSlicePtr := reflect.ValueOf(slicePtr)
	if vSlicePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("slicePtr should be a slice pointer but is %s", vSlicePtr.Type())
	}

	vSlice := vSlicePtr.Elem()
	if vSlice.Kind() != reflect.Slice {
		return fmt.Errorf("slicePtr should be a slice pointer but is %s", vSlicePtr.Type())
	}

	matcher := func(obj interface{}) (bool, error) {
		return p.evaluate(rules, obj)
	}
	return p.filterSliceValue(vSlice, matcher)
}

// filterSliceValue filters a slice passed as a reflect.Value, in place. It calls the matcher function to evaluate whether to keep each item or not.
func (p *processor) filterSliceValue(slice reflect.Value, matcher func(interface{}) (bool, error)) error {
	n := 0

	for i := 0; i < slice.Len(); i++ {
		value := slice.Index(i)

		if !value.CanInterface() {
			return fmt.Errorf("elements contained a %s which cannot be interfaced or set", value.Type())
		}

		match, err := matcher(value.Interface())
		if err != nil {
			return err
		}

		if match {
			// replace unselected elements by the ones that match
			slice.Index(n).Set(value)
			n++
		}
	}

	// shorten the slice to the actual number of elements
	slice.SetLen(n)

	return nil
}

func (p *processor) evaluate(rules [][]Rule, obj interface{}) (bool, error) {
	for i := range rules {
		orResult := false

		for j := range rules[i] {
			// need a pointer to always use the same value and have some state (e.g. regexp)
			rule := &rules[i][j]

			value, err := p.fields.GetFieldValue(rule.Field, obj)
			if err != nil {
				return false, err
			}

			match, err := rule.Evaluate(value)
			if err != nil {
				return false, err
			}

			if match {
				orResult = true
				break
			}
		}

		if !orResult {
			return false, nil
		}
	}

	return true, nil
}
