package errors

import (
	"google.golang.org/grpc/codes"
)

type invalidArgument struct {
	*richError
}

func NewInvalidArgument(message string) RichError {
	return &invalidArgument{
		richError: newRichError(message),
	}
}

func (i *invalidArgument) Error() string {
	return i.message
}

func (i *invalidArgument) ToGRPCError() error {
	return i.toGRPCErrorWith(codes.InvalidArgument)
}
