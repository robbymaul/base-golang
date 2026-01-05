package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ErrorTrace struct {
	Err        error  `json:"-"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func NewErrorTrace(err error) *ErrorTrace {
	return &ErrorTrace{
		Err:        err,
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}

func (e *ErrorTrace) Error() string {
	return e.Message
}

func (e *ErrorTrace) Unwrap() error {
	return e.Err
}

func (e *ErrorTrace) SetStatusCode(code int) *ErrorTrace {
	e.StatusCode = code
	return e
}

func HandleError(c *gin.Context, err error) {
	log.Error().Err(err).Msg(err.Error())

	var trace *ErrorTrace
	if errors.As(err, &trace) {
		log.Error().Err(err).Int("http-status", trace.StatusCode).Msg(trace.Message)
		c.JSON(trace.StatusCode, gin.H{
			"success": false,
			"error":   trace,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error": gin.H{
			"message": "Internal server error",
		},
	})
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			HandleError(c, c.Errors.Last().Err)
			c.Abort()
		}
	}
}
