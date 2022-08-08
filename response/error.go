package response

import "net/http"

type ApiError struct {
	Error  string `json:"error"`
	Code   string `json:"code"`
	Status int    `json:"status"`
}

var (
	NotFoundError           = ApiError{Error: "Resource not found", Code: "NOT_FOUND_ERROR", Status: http.StatusNotFound}
	InternalServerError     = ApiError{Error: "Internal server error", Code: "INTERNAL_SERVER_ERROR", Status: http.StatusInternalServerError}
	BadRequestError         = ApiError{Error: "Invalid parse user input", Code: "BAD_REQUEST", Status: http.StatusBadRequest}
	EmailAlreadyInUse       = ApiError{Error: "Email already in use", Code: "EMAIL_IN_USE", Status: http.StatusConflict}
	ResourceNotFoundError   = ApiError{Error: "Resource not found", Code: "RESOURCE_NOT_FOUND", Status: http.StatusNotFound}
	InvalidCredentialsError = ApiError{Error: "Invalid credentials", Code: "INVALID_CREDENTIALS", Status: http.StatusBadRequest}
	InvalidTokenError       = ApiError{Error: "Invalid token", Code: "INVALID_TOKEN", Status: http.StatusUnauthorized}
)
