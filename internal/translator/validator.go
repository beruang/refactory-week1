package translator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func registerTranslator() (*validator.Validate, error) {
	validate := validator.New()

	english := en.New()
	translator = ut.New(english, english)

	translatorEng, ok := translator.GetTranslator("en")
	if !ok {
		return nil, errors.New("cannot found message translator")
	}

	if err := enTranslation.RegisterDefaultTranslations(validate, translatorEng); nil != err {
		return nil, errors.Wrap(err, "registering default translation")
	}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate, nil
}

var translator *ut.UniversalTranslator

var validate, err = registerTranslator()

func GetValidator() (*validator.Validate, error) {
	return validate, err
}

func GetTranslator() *ut.UniversalTranslator {
	return translator
}
