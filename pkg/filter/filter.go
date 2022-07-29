// Package filter provides generic filtering capabilities for struct slices.
package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
)

const (
	defaultMaxRules = 3

	// DefaultURLQueryFilterKey is the default URL query key used by Processor.ParseURLQuery().
	// Can be customized with WithQueryFilterKey().
	DefaultURLQueryFilterKey = "filter"

	// defaultMaxResults is the default number of results for Apply. Can be overridden with WithMaxResults().
	defaultMaxResults = math.MaxInt
)

// ParseJSON parses and returns a [][]Rule from its JSON representation.
func ParseJSON(s string) ([][]Rule, error) {
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
func New(opts ...Option) (*Processor, error) {
	p := &Processor{
		maxRules:          defaultMaxRules,
		maxResults:        defaultMaxResults,
		urlQueryFilterKey: DefaultURLQueryFilterKey,
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// Processor provides the filtering logic and methods.
type Processor struct {
	fields            fieldGetter
	maxRules          int
	maxResults        int
	urlQueryFilterKey string
}

// ParseURLQuery parses and returns the defined query parameter from a *url.URL.
// Defaults to DefaultURLQueryFilterKey and can be customized with WithQueryFilterKey().
//
// If the query parameter is empty or missing, will return a nil slice.
// If there is a value which is invalid, will return an error.
func (p *Processor) ParseURLQuery(q url.Values) ([][]Rule, error) {
	value := q.Get(p.urlQueryFilterKey)
	if value == "" {
		return nil, nil
	}

	return ParseJSON(value)
}

// Apply filters the slice to remove elements not matching the defined rules.
// The slice parameter must be a pointer to a slice and is filtered *in place*.
//
// This is a shortcut to ApplySubset with 0 offset and maxResults length.
func (p *Processor) Apply(rules [][]Rule, slicePtr interface{}) error {
	return p.ApplySubset(rules, slicePtr, 0, p.maxResults)
}

// ApplySubset filters the slice to remove elements not matching the defined rules.
// The slice parameter must be a pointer to a slice and is filtered *in place*.
//
// Depending on offset, the first results are filtered even if they match
// Depending on length, the filtered slice will only contain a set number of elements.
func (p *Processor) ApplySubset(rules [][]Rule, slicePtr interface{}, offset, length int) error {
	if offset < 0 {
		return errors.New("offset must be positive")
	}

	if length < 1 {
		return errors.New("length must be strictly positive")
	}

	if len(rules) == 0 {
		return nil
	}

	err := p.checkRulesCount(rules)
	if err != nil {
		return err
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
		return p.evaluateRules(rules, obj)
	}

	return p.filterSliceValue(vSlice, offset, length, matcher)
}

func (p *Processor) checkRulesCount(rules [][]Rule) error {
	count := 0
	for i := range rules {
		count += len(rules[i])
	}

	if count > p.maxRules {
		return fmt.Errorf("too many rules: got %d max is %d", count, p.maxRules)
	}

	return nil
}

// filterSliceValue filters a slice passed as a reflect.Value, in place.
// It calls the matcher function to evaluate whether to keep each item or not.
func (p *Processor) filterSliceValue(slice reflect.Value, offset, length int, matcher func(interface{}) (bool, error)) error {
	// number of matched elements
	n := 0

	sliceLen := slice.Len()
	if offset > sliceLen-1 { // offset is out of bounds
		slice.SetLen(0)
		return nil
	}

	for i := 0; i < sliceLen && n < length; i++ {
		value := slice.Index(i)

		// value can always be Interface() because it's in a slice and cannot point to an unexported field
		match, err := matcher(value.Interface())
		if err != nil {
			return err
		}

		if match {
			if offset == 0 {
				// replace unselected elements by the ones that match
				slice.Index(n).Set(value)
				n++
			}

			if offset > 0 {
				offset--
			}
		}
	}

	// shorten the slice to the actual number of elements
	slice.SetLen(n)

	return nil
}

// nolint: gocognit
func (p *Processor) evaluateRules(rules [][]Rule, obj interface{}) (bool, error) {
	for i := range rules {
		orResult := false

		for j := range rules[i] {
			match, err := p.evaluateRule(&rules[i][j], obj)
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

// evaluateRule evaluates a specific rule over an object.
//
// It needs a pointer to let the Rule reuse its state (e.g. precompiled regexp).
func (p *Processor) evaluateRule(rule *Rule, obj interface{}) (bool, error) {
	value, err := p.fields.GetFieldValue(obj, rule.Field)
	if errors.Is(err, errFieldNotFound) {
		return false, nil // filter out missing field without error
	}

	if err != nil {
		return false, err
	}

	return rule.Evaluate(value)
}
