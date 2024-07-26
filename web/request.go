package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Decoder interface {
	Decode(data []byte) error
}

type validator interface {
	Validate() error
}

func Decode(r *http.Request, d Decoder) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request: failed to read request body: %w", err)
	}

	err = d.Decode(data)
	if err != nil {
		return fmt.Errorf("request: failed to decode request body: %w", err)
	}

	if v, ok := d.(validator); ok {
		err = v.Validate()
		if err != nil {
			return fmt.Errorf("request: failed to validate request body: %w", err)
		}
	}

	return nil
}

func DecodeJSONErrorToResponse(reqID string, err error) Error {
	if err == nil {
		return NewInternalServerErrorResponse(reqID)
	}
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return NewBadRequestErrorResponse(reqID, NewFieldsError("body", "Invalid body"))
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return NewBadRequestErrorResponse(reqID, NewFieldsError("body", "Invalid body"))
	}
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return NewBadRequestErrorResponse(reqID, NewFieldsError(
			unmarshalTypeError.Field,
			"Invalid value type for field",
		))
	}
	if errors.Is(err, io.EOF) {
		return NewBadRequestErrorResponse(reqID, NewFieldsError("body", "Empty body"))
	}
	return NewInternalServerErrorResponse(reqID)
}
