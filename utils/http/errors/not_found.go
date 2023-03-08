package errors

import "fmt"

type NotFound struct {
	message string
}

func NewNotFound(message string) *NotFound {
	return &NotFound{
		message: message,
	}
}

func (n *NotFound) Error() string {
	return fmt.Sprintf("HTTP not found error with message: %s", n.message)
}
