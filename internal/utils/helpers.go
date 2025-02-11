package helpers

import "fmt"

// WrapError wraps an error with additional context.
func WrapError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
