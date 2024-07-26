package web

import "net/http"

// 400
func NewBadRequestErrorResponse(requestID string, details FieldErrors) ErrorResponse {
	return ErrorResponse{
		ID:     requestID,
		Status: http.StatusBadRequest,
		Err: ErrorDetail{
			Code:    "400",
			Message: "Bad Request",
			Details: details,
		},
	}
}

// 404
func NewNotFoundErrorResponse(requestID string) ErrorResponse {
	return ErrorResponse{
		ID:     requestID,
		Status: http.StatusNotFound,
		Err: ErrorDetail{
			Code:    "404",
			Message: "Not Found",
		},
	}
}

// 422
func NewUnprocessableEntityResponse(id string, code string) ErrorResponse {
	return ErrorResponse{
		ID:     id,
		Status: http.StatusUnprocessableEntity,
		Err: ErrorDetail{
			Code:    code,
			Message: "Unprocessable Entity",
		},
	}
}

// 500
func NewInternalServerErrorResponse(id string) ErrorResponse {
	return ErrorResponse{
		ID:     id,
		Status: http.StatusInternalServerError,
		Err: ErrorDetail{
			Code:    "500",
			Message: "Internal Server Error",
		},
	}
}
