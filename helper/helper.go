package helper

import (
	"github.com/go-playground/validator/v10"
)

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func ResponseAPI(message string, code int, status string, data interface{}) (resp Response) {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	resp = Response{
		Meta: meta,
		Data: data,
	}
	return
}

func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
