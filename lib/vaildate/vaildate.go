package vaildate

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

// field validator for post method
var bindingValidator *validator.Validate

// field validator for patch or put method
var updateValidator *validator.Validate

// validate err translator
var BindingTrans ut.Translator
var UpdateTrans ut.Translator

func Init(locale string) {
	bindingValidator, _ = binding.Validator.Engine().(*validator.Validate)
	en0 := en.New()
	zh0 := zh.New()
	uni0 := ut.New(en0, zh0)
	//預設語系
	BindingTrans, _ = uni0.GetTranslator(locale)
	switch locale {
	case "zh":
		zh_translations.RegisterDefaultTranslations(bindingValidator, BindingTrans)
		break
	case "en":
		en_translations.RegisterDefaultTranslations(bindingValidator, BindingTrans)
		break
	default:
		zh_translations.RegisterDefaultTranslations(bindingValidator, BindingTrans)
		break
	}
	bindingValidator.RegisterTagNameFunc(registerTagName)

	updateValidator = validator.New()
	updateValidator.SetTagName("update")
	en1 := en.New()
	zh1 := zh.New()
	uni1 := ut.New(en1, zh1)
	// register err message English
	UpdateTrans, _ = uni1.GetTranslator(locale)
	switch locale {
	case "zh":
		zh_translations.RegisterDefaultTranslations(updateValidator, UpdateTrans)
		break
	case "en":
		en_translations.RegisterDefaultTranslations(updateValidator, UpdateTrans)
		break
	default:
		zh_translations.RegisterDefaultTranslations(updateValidator, UpdateTrans)
		break
	}
	updateValidator.RegisterTagNameFunc(registerTagName)
	updateValidator.RegisterValidation("fixed", fixed)
	updateValidator.RegisterTranslation("fixed", UpdateTrans, fixedTranslation, fixedTranslationAdding)
}

func registerTagName(field reflect.StructField) string {
	name := field.Tag.Get("json")
	fmt.Println(name)
	if name == "-" {
		return ""
	}
	return name
}

//this field must be read only, can not be update
func fixed(f1 validator.FieldLevel) bool {
	return false
}

// fixed tag custom error message
func fixedTranslation(ut ut.Translator) error {
	return ut.Add("fixed", "{0} can not be updated. it's fixed.", true)
}

// register fixed tag custom error message
func fixedTranslationAdding(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("fixed", fe.Field())
	return t
}
