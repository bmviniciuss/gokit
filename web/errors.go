package web

import (
	"encoding/json"
	"errors"
)

type ErrorResponse struct {
	Status int         `json:"-"`
	ID     string      `json:"id"`
	Err    ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details FieldErrors `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrors []FieldError

func NewFieldsError(field string, message string) FieldErrors {
	return FieldErrors{
		{
			Field:   field,
			Message: message,
		},
	}
}

func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}
