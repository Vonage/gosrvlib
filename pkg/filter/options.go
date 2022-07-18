package filter

import "errors"

// Option is the function that allows to set configuration options.
type Option func(v *processor) error

// WithFieldNameTag allows to use the field names specified by the tag instead of the original struct names.
func WithFieldNameTag(tag string) Option {
	return func(v *processor) error {
		if tag == "" {
			return errors.New("tag cannot be empty")
		}

		v.fields.fieldTag = tag

		return nil
	}
}
