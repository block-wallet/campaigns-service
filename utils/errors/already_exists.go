package errors

import (
	"google.golang.org/grpc/codes"
)

type alreadyExists struct {
	*richError
}

func NewAlreadyExists(message string) RichError {
	return &alreadyExists{
		newRichError(message),
	}
}

func (a *alreadyExists) Error() string {
	return a.message
}

func (a *alreadyExists) ToGRPCError() error {
	return a.toGRPCErrorWith(codes.AlreadyExists)
}
