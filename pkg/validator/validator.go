package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidator() (*Validator, error) {
	uni := ut.New(en.New(), en.New())
	v := &Validator{}

	trans, _ := uni.GetTranslator("en")
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, fmt.Errorf("failed to translations.RegisterDefaultTranslations: %w", err)
	}

	v.validator = validate
	v.translator = trans

	return v, nil
}

type CustomValidationError struct {
	Message map[string]string `json:"message"`
}

func NewCustomValidationError() *CustomValidationError {
	return &CustomValidationError{
		Message: make(map[string]string),
	}
}

func (e CustomValidationError) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func (v *Validator) Validate(i any) (err error) {
	if err = v.validator.Struct(i); err != nil {
		var vErrs validator.ValidationErrors

		ve := NewCustomValidationError()
		if errors.As(err, &vErrs) {
			for _, e := range vErrs {
				ve.Message[strings.ToLower(e.Field())] = e.Translate(v.translator)
			}
			return ve
		}

		return err
	}

	return nil
}
