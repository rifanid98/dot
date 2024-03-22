package util

import (
	"regexp"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

// NewValidator Initiatilize validator in singleton way
func NewValidator() *validator.Validate {
	if validate != nil {
		return validate
	}
	en := en.New()
	uni := ut.New(en, en)

	translator, _ = uni.GetTranslator("en")
	validate = validator.New()

	registerValidation()
	registerTranslation()

	err := en_translations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		panic(err)
	}

	return validate
}

func registerValidation() {
	if err := validate.RegisterValidation("ISO8601", func(fl validator.FieldLevel) bool {
		ISO8601RegexString := "^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)$"
		ISO8601Regex := regexp.MustCompile(ISO8601RegexString)
		return ISO8601Regex.MatchString(fl.Field().String())
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("crmpwd", func(fl validator.FieldLevel) bool {
		ISO8601RegexString := "^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)$"
		ISO8601Regex := regexp.MustCompile(ISO8601RegexString)
		return ISO8601Regex.MatchString(fl.Field().String())
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("majoolitepwd", func(fl validator.FieldLevel) bool {
		ISO8601RegexString := "^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)$"
		ISO8601Regex := regexp.MustCompile(ISO8601RegexString)
		return ISO8601Regex.MatchString(fl.Field().String())
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("dashboardpwd", func(fl validator.FieldLevel) bool {
		ISO8601RegexString := "^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)$"
		ISO8601Regex := regexp.MustCompile(ISO8601RegexString)
		return ISO8601Regex.MatchString(fl.Field().String())
	}); err != nil {
		panic(err)
	}
}

func registerTranslation() {
	if err := validate.RegisterTranslation("ISO8601", translator, func(ut ut.Translator) error {
		return ut.Add("ISO8601", "{0} must following ISO8601 or RFC3339 date format", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ISO8601", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterTranslation("crmpwd", translator, func(ut ut.Translator) error {
		return ut.Add("ISO8601", "{0} must following ISO8601 or RFC3339 date format", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ISO8601", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterTranslation("majoolitepwd", translator, func(ut ut.Translator) error {
		return ut.Add("ISO8601", "{0} must following ISO8601 or RFC3339 date format", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ISO8601", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterTranslation("dashboardpwd", translator, func(ut ut.Translator) error {
		return ut.Add("ISO8601", "{0} must following ISO8601 or RFC3339 date format", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ISO8601", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}
}

func GetValidatorMessage(err error) (messages []string) {
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, err.Translate(translator))
		}
	}
	return
}
