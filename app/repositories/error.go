package repositories

import (
	"fmt"
	"strings"
)

type Error struct {
	Message string
	Reason  string
}

func (e *Error) Error() string {
	var messages []string
	if e.Reason != "" {
		messages = append(messages, e.Reason)
	}

	if len(messages) > 0 {
		return fmt.Sprintf("repository: %s (%s)", e.Message, strings.Join(messages, "; "))
	}

	return fmt.Sprintf("repository: %s", e.Message)
}

func newError(context string, reason string) error {
	return &Error{
		Message: context,
		Reason:  reason,
	}
}
