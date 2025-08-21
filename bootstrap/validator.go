package bootstrap

import (
	"context"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vi_translations "github.com/go-playground/validator/v10/translations/vi"
)

const (
	VI           = "vi"
	regexPhoneVi = `^(84|0(3|5|7|8|9))[0-9]{8}$`
)

type ValidationError struct {
	Message string
	Data    map[string]any
}

func (e *ValidationError) Error() string {
	return e.Message
}

type IValidator interface {
	IsPhoneNumber(local, val string) bool
	ValidateStruct(s any) *ValidationError
	EmailOrPhone(local, val string) bool
}

type customValidator struct {
	v     *validator.Validate
	trans ut.Translator
}

func RegisterCustomValidations(v *validator.Validate) IValidator {
	viLocale := vi.New()
	uni := ut.New(viLocale, viLocale)
	trans, _ := uni.GetTranslator("vi")
	vi_translations.RegisterDefaultTranslations(v, trans)
	cv := &customValidator{v, trans}
	// -----
	v.RegisterValidation("tell", func(fl validator.FieldLevel) bool { return cv.IsPhoneNumber(fl.Param(), fl.Field().String()) })
	v.RegisterTranslation("tell", cv.trans, func(ut ut.Translator) error {
		return ut.Add("tell", "Trường này phải là số điện thoại hợp lệ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		msg, _ := ut.T("tell", fe.Field())
		return msg
	})

	v.RegisterValidation("email_or_tell", func(fl validator.FieldLevel) bool { return cv.EmailOrPhone(fl.Param(), fl.Field().String()) })
	v.RegisterTranslation("email_or_tell", cv.trans, func(ut ut.Translator) error {
		return ut.Add("email_or_tell", "Trường này phải là email hoặc số điện thoại hợp lệ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		msg, _ := ut.T("email_or_tell", fe.Field())
		return msg
	})

	v.RegisterValidation("in", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		param := fl.Param()
		params := strings.Fields(param)
		return slices.Contains(params, val)
	})
	v.RegisterTranslation("in", cv.trans, func(ut ut.Translator) error {
		return ut.Add("in", "Trường này phải là một trong các giá trị: {0}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		msg, _ := ut.T("in", fe.Param())
		return msg
	})

	return cv
}

func (v *customValidator) IsPhoneNumber(local, val string) bool {
	switch local {
	case VI:
		regexp := regexp.MustCompile(regexPhoneVi)
		return regexp.MatchString(val)
	default:
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{9,15}$`)
		return phoneRegex.MatchString(val)
	}
}

func (v *customValidator) EmailOrPhone(local, val string) bool {
	v.v.VarCtx(context.Background(), val, "email")
	if err := v.v.VarCtx(context.Background(), val, "email"); err == nil {
		return true
	}
	return v.IsPhoneNumber(local, val)
}

func (v *customValidator) ValidateStruct(s any) *ValidationError {
	err := v.v.Struct(s)
	if err != nil {
		var message string
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, ve := range validationErrors {
			field := ve.Field()
			errorMessages[field] = ve.Translate(v.trans)
		}
		data := make(map[string]any, len(errorMessages))
		for k, v := range errorMessages {
			data[k] = v
			message += k + ": " + v + "; "
		}

		return &ValidationError{
			Message: message,
			Data:    data,
		}
	}
	return nil
}
