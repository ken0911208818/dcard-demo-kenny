package vaildate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
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
	uni0 := ut.New(en0, zh.New())
	//預設語系
	BindingTrans, _ = uni0.GetTranslator(locale)
	enTranslations.RegisterDefaultTranslations(bindingValidator, BindingTrans)
	bindingValidator.RegisterTagNameFunc(registerTagName)

}

func registerTagName(field reflect.StructField) string {
	name := field.Tag.Get("json")
	if name == "-" {
		return ""
	}
	return name
}
