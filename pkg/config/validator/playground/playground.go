package playground

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type (
	Validator struct {
		validate *validator.Validate
		trans    ut.Translator
	}
)

func New(local string) *Validator {
	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT, enT)

	if local == "" {
		local = "zh"
	}

	trans, _ := uni.GetTranslator(local)
	validate := validator.New()

	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		validate: validate,
		trans:    trans,
	}
}

func (v Validator) Validate(val interface{}) error {
	err := v.validate.Struct(val)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)

	var errMsg []string
	for _, msg := range errs.Translate(v.trans) {
		errMsg = append(errMsg, msg)
	}

	return errors.New(strings.Join(errMsg, ";"))
}
