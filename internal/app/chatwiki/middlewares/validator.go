// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package middlewares

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/zhimaAi/go_tools/logs"
	"go.uber.org/zap"
)

var (
	uni *ut.UniversalTranslator
	//trans ut.Translator
)

func Init() {
	//register editor
	zhT := zh.New()
	enT := en.New()
	uni = ut.New(zhT, enT)
}
func switchTrans(lang string) ut.Translator {
	var err error
	trans, _ := uni.GetTranslator(lang)

	//get gin validator
	validate := binding.Validator.Engine().(*validator.Validate)
	// check trans
	switch lang {
	case "en":
		err = entranslations.RegisterDefaultTranslations(validate, trans)
	case define.LangZhCn:
		err = zhtranslations.RegisterDefaultTranslations(validate, trans)
	default:
		err = zhtranslations.RegisterDefaultTranslations(validate, trans)
	}
	if err != nil {
		logs.Error("validator register error", zap.Error(err))
	}
	return trans
}

func GetValidateErr(obj any, rawErr error, lang string) error {
	var validationErrs validator.ValidationErrors
	if !errors.As(rawErr, &validationErrs) {
		return rawErr
	}
	trans := switchTrans(lang)
	var errString []string
	for _, validationErr := range validationErrs {
		field, ok := reflect.TypeOf(obj).FieldByName(validationErr.Field())
		if ok {
			if e := field.Tag.Get("msg"); e != "" {
				errString = append(errString, fmt.Sprintf("%s: %s", field.Name, e))
				continue
			} else {
				errString = append(errString, fmt.Sprintf("%s", validationErr.Translate(trans)))
			}
		} else {
			errString = append(errString, validationErr.Error())
		}
	}
	return errors.New(strings.Join(errString, ";"))
}
