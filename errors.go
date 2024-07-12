package fast

import (
	"fmt"
	"net/http"
)

// httpError is an error that contains an HTTP status code and a message.
type httpError struct {
	status  int
	message string
}

func (e httpError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.status, e.message)
}

func ValidationError(message string) httpError {
	return httpError{
		status:  http.StatusUnprocessableEntity,
		message: message,
	}
}

func UnauthorizedError(message string) httpError {
	return httpError{
		status:  http.StatusUnauthorized,
		message: message,
	}
}
