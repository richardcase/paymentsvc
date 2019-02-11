package model

// ErrorResponse represents an error message response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
