package client

import "fmt"

// APIError represents an error returned by the PixiGPT API (OpenAI format).
type APIError struct {
	ErrorData struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
	StatusCode int `json:"-"` // HTTP status code
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.ErrorData.Code != "" {
		return fmt.Sprintf("[%s] %s: %s", e.ErrorData.Code, e.ErrorData.Type, e.ErrorData.Message)
	}
	return fmt.Sprintf("[%s] %s", e.ErrorData.Type, e.ErrorData.Message)
}

// IsAuthError returns true if the error is an authentication error.
func IsAuthError(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.ErrorData.Type == "authentication_error"
}

// IsRateLimitError returns true if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.ErrorData.Type == "rate_limit_error"
}
