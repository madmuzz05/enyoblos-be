package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = func() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}()

func ValidateRequest(ctx fiber.Ctx, req interface{}) (map[string]string, error) {
	if err := ctx.Bind().Body(req); err != nil {
		return nil, err
	}

	if err := validate.Struct(req); err != nil {
		return FormatValidationError(err), err
	}

	return nil, nil
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()

			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s wajib diisi", field)
			case "email":
				errors[field] = fmt.Sprintf("%s tidak valid", field)
			case "min":
				errors[field] = fmt.Sprintf("minimal karakter %s", e.Param())
			case "max":
				errors[field] = fmt.Sprintf("melebihi batas maksimal karakter %s", e.Param())
			default:
				errors[field] = fmt.Sprintf("%s tidak valid", field)
			}
		}
	}

	return errors
}
