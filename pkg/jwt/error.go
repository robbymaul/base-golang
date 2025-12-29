package pkgjwt

import "fmt"

func Error(err error) error {
	return fmt.Errorf("failed to verify JWT token: %w", err)
}
