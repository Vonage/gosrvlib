// Package validator expose wrapper function for https://github.com/go-playground/validator
// to provide value validations for structs and individual fields based on tags.
package validator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	ut "github.com/go-playground/universal-translator"
	vt "github.com/go-playground/validator/v10"
	"go.uber.org/multierr"
)

// ValidationError is a custom error adding a Field member.
type ValidationError struct {
	// Tag is the validation tag that failed.
	// If the validation was an alias, this will return the alias name and not the underlying tag that failed.
	//
	// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
	// will return "iscolor"
	Tag string

	// ActualTag is the validation tag that failed,
	// even if an alias the actual tag within the alias will be returned.
	// If an 'or' validation fails the entire or will be returned.
	//
	// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
	// will return "hexcolor|rgb|rgba|hsl|hsla"
	ActualTag string

	// Namespace for the field error,
	// with the tag name taking precedence over the field's actual name.
	Namespace string

	// StructNamespace is the namespace for the field error,
	// with the field's actual name.
	StructNamespace string

	// Field is the field name with the tag name taking precedence over the field's actual name.
	Field string

	// StructField is the field's actual name from the struct, when able to determine.
	StructField string

	// Value the actual field's value
	Value interface{}

	// Param is the param value
	Param string

	// Kind returns the Field's reflect Kind in string format
	Kind string

	// Error returns the translated error message
	Err string
}

// Error returns a string representation of the error.
func (e *ValidationError) Error() string {
	return e.Err
}

// Validator contains the validator object fields.
type Validator struct {
	// V is the validate object
	V *vt.Validate

	// Trans is the translator object
	T ut.Translator

	// translate contains the map of basic translation templates indexed by tag
	translate map[string]*template.Template
}

// New returns a new validator with the specified options.
func New(opts ...Option) (*Validator, error) {
	v := &Validator{
		V: vt.New(),
	}
	for _, applyOpt := range opts {
		if err := applyOpt(v); err != nil {
			return nil, err
		}
	}
	return v, nil
}

// ValidateStruct validates the structure fields tagged with "validate".
func (v *Validator) ValidateStruct(obj interface{}) error {
	err := v.V.Struct(obj)
	if err == nil {
		return nil
	}
	for _, e := range err.(vt.ValidationErrors) {
		if e != nil {
			ve := &ValidationError{
				Tag:             e.Tag(),
				ActualTag:       e.ActualTag(),
				Namespace:       e.Namespace(),
				StructNamespace: e.StructNamespace(),
				Field:           e.Field(),
				StructField:     e.StructField(),
				Value:           e.Value(),
				Param:           e.Param(),
				Kind:            e.Kind().String(),
			}
			ve.Err = v.stringify(e, ve)
			err = multierr.Append(err, ve)
		}
	}
	return err
}

func (v *Validator) stringify(fe vt.FieldError, ve *ValidationError) string {
	if v.T != nil {
		return fe.Translate(v.T)
	}
	if v.translate != nil {
		ns := fe.Namespace()
		if idx := strings.Index(ns, "."); idx != -1 {
			ns = ns[idx+1:] // remove root struct name
		}
		ve.Namespace = ns
		t, ok := v.translate[ve.Tag]
		if ok {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, fe); err == nil {
				return tpl.String()
			}
		}
	}
	return fmt.Sprintf("%s", fe)
}
