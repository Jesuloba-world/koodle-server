package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"

)

func GenericValidate(input interface{}) error {
	validationErr := Validate.Struct(input)
	if validationErr != nil {
		errs := validationErr.(validator.ValidationErrors)

		var errMsg strings.Builder
		for _, e := range errs {
			errMsg.WriteString(e.Translate(Translator))
			errMsg.WriteString("; ")
		}

		return errors.New(errMsg.String())
	}

	return nil
}
