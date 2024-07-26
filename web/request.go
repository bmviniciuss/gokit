package web

import (
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
