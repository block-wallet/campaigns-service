package errors

import (
	"google.golang.org/grpc/codes"
)

type NotFound struct {
	*richError
}

func NewNotFound(message string) RichError {
	return &NotFound{
		richError: newRichError(message),
	}
}

func (n *NotFound) Error() string {
	return n.message
}

func (n *NotFound) ToGRPCError() error {
	return n.toGRPCErrorWith(codes.NotFound)
}
