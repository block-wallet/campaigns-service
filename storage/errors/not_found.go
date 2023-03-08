package errors

import "fmt"

type NotFound struct {
	Key string
}

func NewNotFound(key string) *NotFound {
	return &NotFound{
		Key: key,
	}
}

func (n *NotFound) Error() string {
	return fmt.Sprintf("Entity not found: %s", n.Key)
}
