package filter

import "errors"

// Option is the function that allows to set configuration options.
type Option func(v *processor) error

// WithFieldNameTag allows to use the field names specified by the tag instead of the original struct names.
//
// Returns an error if the tag is empty.
func WithFieldNameTag(tag string) Option {
	return func(v *processor) error {
		if tag == "" {
			return errors.New("tag cannot be empty")
		}

		v.fields.fieldTag = tag

		return nil
	}
}

// WithMaxRules sets the maximum number of rules to pass to the Processor.Apply() function without errors.
// If this option is not set, it defaults to 3.
//
// Return an error if max is less than 1.
func WithMaxRules(max int) Option {
	return func(v *processor) error {
		if max < 1 {
			return errors.New("max")
		}

		v.maxRules = max

		return nil
	}
}
