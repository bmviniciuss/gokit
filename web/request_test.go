package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("forced read error")
}

type decodeStruct struct {
}

func (d *decodeStruct) Decode(data []byte) error {
	return fmt.Errorf("forced decode error")
}

type nonValidable struct {
	Name string `json:"name"`
}

func (d *nonValidable) Decode(data []byte) error {
	return json.Unmarshal(data, d)
}

type validable struct {
	Name string `json:"name"`
}

func (d *validable) Decode(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *validable) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func TestDecode(t *testing.T) {
	t.Run("should return an error if read all fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", errorReader{})
		got := Decode(req, nil)
		assert.Error(t, got)
	})

	t.Run("should return an error Decode fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{}))
		got := Decode(req, &decodeStruct{})
		assert.Error(t, got)
	})

	t.Run("should return nil on decode success and struct does not have validation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"test"}`)))
		var body nonValidable
		got := Decode(req, &body)
		assert.Nil(t, got)
		assert.Equal(t, "test", body.Name)
	})

	t.Run("should return nil after validation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"test"}`)))
		var body validable
		got := Decode(req, &body)
		assert.Nil(t, got)
		assert.Equal(t, "test", body.Name)
	})

	t.Run("should return validation validation error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":""}`)))
		var body validable
		got := Decode(req, &body)
		assert.Error(t, got)
	})
}

func TestDecodeJSONErrorToResponse(t *testing.T) {
	var (
		reqID = uuid.NewString()
	)
	t.Run("should return 500 response on nil error", func(t *testing.T) {
		got := DecodeJSONErrorToResponse(reqID, nil)
		assert.Equal(t, http.StatusInternalServerError, got.Status)
		assert.Equal(t, reqID, got.ID)
		assert.Equal(t, "500", got.Err.Code)
		assert.Equal(t, "Internal Server Error", got.Err.Message)
	})

	t.Run("should return bad request if json is mal formed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{`)))
		var body nonValidable
		err := Decode(req, &body)
		assert.Error(t, err)
		got := DecodeJSONErrorToResponse(reqID, err)
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, reqID, got.ID)
		assert.Equal(t, "400", got.Err.Code)
		assert.Equal(t, "Bad Request", got.Err.Message)
		assert.Equal(t, "body", got.Err.Details[0].Field)
		assert.Equal(t, "Invalid body", got.Err.Details[0].Message)
	})

	t.Run("should return bad request field has different type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":1}`)))
		var body nonValidable
		err := Decode(req, &body)
		assert.Error(t, err)
		got := DecodeJSONErrorToResponse(reqID, err)
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, reqID, got.ID)
		assert.Equal(t, "400", got.Err.Code)
		assert.Equal(t, "Bad Request", got.Err.Message)
		assert.Equal(t, "name", got.Err.Details[0].Field)
		assert.Equal(t, "Invalid value type for field", got.Err.Details[0].Message)
	})

	t.Run("should return bad request if json is empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(``)))
		var body nonValidable
		err := Decode(req, &body)
		assert.Error(t, err)
		got := DecodeJSONErrorToResponse(reqID, err)
		assert.Equal(t, http.StatusBadRequest, got.Status)
		assert.Equal(t, reqID, got.ID)
		assert.Equal(t, "400", got.Err.Code)
		assert.Equal(t, "Bad Request", got.Err.Message)
		assert.Equal(t, "body", got.Err.Details[0].Field)
		assert.Equal(t, "Invalid body", got.Err.Details[0].Message)
	})
}
