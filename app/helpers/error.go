package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/web"
)

// ErrorTrace represents a detailed error trace with context and status code
type ErrorTrace struct {
	Err        error  // Original error
	Message    string // Human-readable message
	Context    string // Where the error occurred
	StatusCode int    // HTTP status code
	Details    any    // Additional error details
}

// NewErrorTrace creates a new ErrorTrace with context
func NewErrorTrace(err error, context string) *ErrorTrace {
	message := err.Error()
	if message == "" {
		message = "unknown error"
	}

	return &ErrorTrace{
		Err:     err,
		Message: message,
		Context: context,
		// Default status code if not set
		StatusCode: http.StatusInternalServerError,
	}
}

// Error implements the error interface
func (e *ErrorTrace) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s: %s", e.Context, e.Message)
	}
	return e.Message
}

// Unwrap returns the underlying error for error inspection
func (e *ErrorTrace) Unwrap() error {
	return e.Err
}

// WithStatusCode sets the HTTP status code and returns the ErrorTrace for chaining
func (e *ErrorTrace) WithStatusCode(code int) *ErrorTrace {
	e.StatusCode = code
	return e
}

// WithDetails adds additional error details and returns the ErrorTrace for chaining
func (e *ErrorTrace) WithDetails(details any) *ErrorTrace {
	e.Details = details
	return e
}

// ErrorResponse handles error responses in a consistent way
func ErrorResponse(ctx *gin.Context, err error) {
	// Log the error with context
	log.Error().Err(err).Msg("request error")

	var trace *ErrorTrace
	if errors.As(err, &trace) {
		response := web.ResponseWeb{
			Success: false,
			Message: trace.Err.Error(),
		}

		// Include details if they exist
		if trace.Details != nil {
			response.Data = trace.Details
		}

		ctx.JSON(trace.StatusCode, response)
		return
	}

	// Fallback for non-ErrorTrace errors
	ctx.JSON(http.StatusInternalServerError, web.ResponseWeb{
		Success: false,
		Message: trace.Err.Error(),
	})
}
