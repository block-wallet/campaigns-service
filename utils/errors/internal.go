package errors

import (
	"google.golang.org/grpc/codes"
)

type Internal struct {
	*richError
}

func NewInternal(message string) RichError {
	return &Internal{
		richError: newRichError(message),
	}
}

func (i *Internal) Error() string {
	return i.message
}

func (i *Internal) ToGRPCError() error {
	return i.toGRPCErrorWith(codes.Internal)
}
