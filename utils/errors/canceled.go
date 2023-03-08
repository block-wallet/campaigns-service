package errors

import (
	"google.golang.org/grpc/codes"
)

type canceled struct {
	*richError
}

func NewCanceled(message string) RichError {
	return &canceled{
		newRichError(message),
	}
}

func (c *canceled) Error() string {
	return c.message
}

func (c *canceled) ToGRPCError() error {
	return c.toGRPCErrorWith(codes.Canceled)
}
