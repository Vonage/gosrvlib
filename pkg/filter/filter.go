// Package filter provides generic filtering capabilities for struct slices
package filter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"github.com/pkg/errors"
)

// Processor is the interface to provide subtractive functions.
type Processor interface {
	// Apply filters the slice to remove elements not matching the defined filters.
	// The slice parameter must be a pointer to a slice and is filtered *in place*.
	Apply(slice interface{}) error
}

// GetFilter returns the "filter" query parameter from a *url.URL.
func GetFilter(u *url.URL) string {
	return u.Query().Get("filter")
}

// ParseRules parses and returns a [][]Rule from its JSON representation.
func ParseRules(s string) ([][]Rule, error) {
	var r [][]Rule
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal json rules")
	}

	return r, nil
}

// New returns a new Processor with the rules and the given options.
//
// The first level of rules is matched with an AND operator and the second level with an OR.
//
// "[a,[b,c],d]" evaluates to "a AND (b OR c) AND d".
func New(r [][]Rule, opts ...Option) (Processor, error) {
	f := processor{
		rules: r,
	}

	for _, opt := range opts {
		if err := opt(&f); err != nil {
			return nil, err
		}
	}

	return &f, nil
}

type processor struct {
	rules  [][]Rule
	fields fieldGetter
}

func (p *processor) Apply(slicePtr interface{}) error {
	if len(p.rules) == 0 {
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

	err := p.filterSliceValue(vSlice, p.evaluateValue)
	if err != nil {
		return err
	}

	return nil
}

// filterSliceValue filters a slice passed as a reflect.Value, in place. It calls the matcher function to evaluate whether to keep each item or not.
func (p *processor) filterSliceValue(slice reflect.Value, matcher func(reflect.Value) (bool, error)) error {
	n := 0

	for i := 0; i < slice.Len(); i++ {
		value := slice.Index(i)

		match, err := matcher(value)
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

func (p *processor) evaluateValue(value reflect.Value) (bool, error) {
	if !value.CanInterface() {
		return false, fmt.Errorf("elements contained a %s which cannot be interfaced or set", value.Type())
	}

	return p.evaluate(value.Interface())
}

func (p *processor) evaluate(obj interface{}) (bool, error) {
	for i := range p.rules {
		orResult := false

		for j := range p.rules[i] {
			// need a pointer to always use the same value and have some state (e.g. regexp)
			rule := &p.rules[i][j]

			value, err := p.fields.getFieldValue(rule.Field, obj)
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
