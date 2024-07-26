package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFieldsError(t *testing.T) {
	tests := []struct {
		field    string
		message  string
		expected FieldErrors
	}{
		{
			field:   "username",
			message: "cannot be empty",
			expected: FieldErrors{
				{
					Field:   "username",
					Message: "cannot be empty",
				},
			},
		},
		{
			field:   "password",
			message: "must be at least 8 characters",
			expected: FieldErrors{
				{
					Field:   "password",
					Message: "must be at least 8 characters",
				},
			},
		},
		{
			field:   "email",
			message: "is invalid",
			expected: FieldErrors{
				{
					Field:   "email",
					Message: "is invalid",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			got := NewFieldsError(tt.field, tt.message)
			assert.Equal(t, tt.expected, got)
		})
	}
}
