package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

type V10 struct {
	validate *v10.Validate
	trans    ut.Translator
}

var (
	ErrFailedToGetTranslator = errors.New("failed to get translator for `en`")
	ErrValidationFailed      = errors.New("validation failed")
	ErrInputNotTranslatable  = errors.New("input is not translatable")
)

func NewValidatorV10() (*V10, error) {
	validate := v10.New()

	english := en.New()
	uni := ut.New(english, english)

	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, ErrFailedToGetTranslator
	}

	if err := entranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0] //nolint:gomnd

		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})

	instance := V10{
		validate: validate,
		trans:    trans,
	}

	return &instance, nil
}

// ValidateStruct implements validator.Validator interface.
func (v V10) ValidateStruct(input any) error {
	if err := v.validate.Struct(input); err != nil {
		return errors.Join(ErrValidationFailed, err)
	}

	return nil
}

// Translate implements validator.Validator interface.
func (v V10) Translate(err error) []Error {
	var verrs v10.ValidationErrors

	if !errors.As(err, &verrs) {
		return []Error{
			{
				Message: err.Error(),
			},
		}
	}

	errs := make([]Error, len(verrs))

	for idx, field := range verrs {
		errs[idx] = Error{
			Field:   field.Field(),
			Message: field.Translate(v.trans),
		}
	}

	return errs
}
