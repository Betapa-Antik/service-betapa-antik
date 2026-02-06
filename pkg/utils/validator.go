package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func NewCustomValidator() echo.Validator {
	v := validator.New()

	// ✅ Register custom password validator
	v.RegisterValidation("password", PasswordValidator)
	return &CustomValidator{
		validator: v,
	}
}

// ParseValidationError converts validator errors to readable messages
func ParseValidationError(err error) map[string]string {
	messages := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()

			// Generate readable message based on tag
			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s wajib diisi", fieldName)
			case "min":
				message = fmt.Sprintf("%s minimal %s karakter", fieldName, fieldError.Param())
			case "max":
				message = fmt.Sprintf("%s maksimal %s karakter", fieldName, fieldError.Param())
			case "email":
				message = fmt.Sprintf("%s harus format email", fieldName)
			case "url":
				message = fmt.Sprintf("%s harus format URL", fieldName)
			case "password":
				message = fmt.Sprintf("%s minimal 8 karakter dan harus ada huruf besar, kecil, angka, dan simbol", fieldName)

			case "eqfield":
				message = fmt.Sprintf("%s harus sama dengan %s", fieldName, fieldError.Param())
			default:
				message = fmt.Sprintf("%s validasi gagal pada %s", fieldName, tag)
			}

			messages[fieldName] = message
		}
	} else {
		// Fallback untuk error non-validator
		messages["error"] = err.Error()
	}

	return messages
}

// // GetValidationErrorMessage returns first validation error message
// func GetValidationErrorMessage(err error) string {
// 	messages := ParseValidationError(err)
// 	for _, msg := range messages {
// 		return msg
// 	}
// 	return "Validasi gagal"
// }
