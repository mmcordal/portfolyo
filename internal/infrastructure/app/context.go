package app

import (
	"portfolyo/internal/infrastructure/errorsx"
	"portfolyo/internal/viewmodel"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Ctx struct {
	*fiber.Ctx
}

func (c *Ctx) BodyParseValidate(m interface{}) []error {
	if err := c.BodyParser(m); err != nil {
		return []error{err}
	}

	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("labelName"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	var errs []error
	err := v.Struct(m)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, e)
		}
		return errs
	}

	vmValidation, ok := m.(viewmodel.Validation)
	if ok {
		errs = append(errs, vmValidation.Validate()...)
	}
	return errs
}

func (c *Ctx) SuccessResponse(data interface{}, dataCount int, message string) errorsx.APIError {
	m := &viewmodel.SuccessResponse{
		Message:   message,
		Success:   true,
		DataCount: dataCount,
		Data:      data,
	}
	err := c.JSON(m)
	if err != nil {
		return errorsx.InternalError(err)
	}

	return nil
}
