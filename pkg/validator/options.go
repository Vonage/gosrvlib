package validator

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	lc "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vt "github.com/go-playground/validator/v10"
	tr "github.com/go-playground/validator/v10/translations/en"
)

// Option is the interface that allows to set options.
type Option func(v *Validator) error

// WithFieldNameTag allows to use the field names specified by the tag instead of the original struct names.
func WithFieldNameTag(tag string) Option {
	return func(v *Validator) error {
		if tag == "" {
			return nil
		}
		v.V.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		return nil
	}
}

// WithBasicTranslations enable the basic internal error message translations.
// The argument t maps tag to the error template using the ValidationError data.
func WithBasicTranslations(t map[string]*template.Template) Option {
	return func(v *Validator) error {
		v.translate = t
		return nil
	}
}

// WithDefaultTranslations sets the default English translations.
func WithDefaultTranslations() Option {
	return func(v *Validator) error {
		en := lc.New()
		uni := ut.New(en, en)
		trans, ok := uni.GetTranslator("en")
		if ok {
			_ = tr.RegisterDefaultTranslations(v.V, trans)
			v.T = trans
		}
		return nil
	}
}

// WithValidationTranslated allows to register a validation func and a translation for the provided tag.
func WithValidationTranslated(tag string, fn vt.Func, registerFn vt.RegisterTranslationsFunc, translationFn vt.TranslationFunc) Option {
	return func(v *Validator) error {
		if err := v.V.RegisterValidation(tag, fn); err != nil {
			return err
		}
		if v.T == nil {
			return fmt.Errorf("the Translator object is nil")
		}
		if err := v.V.RegisterTranslation(tag, v.T, registerFn, translationFn); err != nil {
			return err
		}
		return nil
	}
}
