package validator

import (
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

)

var Validate *validator.Validate
var Translator ut.Translator

func init() {
	en := en.New()
	uni := ut.New(en, en)
	Translator, _ = uni.GetTranslator("en")

	Validate = validator.New()
	en_translations.RegisterDefaultTranslations(Validate, Translator)

	// custom function for passsword validation
	Validate.RegisterValidation("password", func(f1 validator.FieldLevel) bool {
		password := f1.Field().String()

		// RULE: a mix of letters, numbers and symbols
		var (
			hasLetter bool
			hasNumber bool
			hasSymbol bool
		)

		for _, char := range password {
			switch {
			case unicode.IsLetter(char):
				hasLetter = true
			case unicode.IsDigit(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSymbol = true
			}

			if hasLetter && hasNumber && hasSymbol {
				return true
			}
		}

		return hasLetter && hasNumber && hasSymbol
	})

	// custom error message for password
	Validate.RegisterTranslation("password", Translator, func(ut ut.Translator) error {
		return ut.Add("password", "{0} must contain a mix of letters, numbers, and symbols", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
}
