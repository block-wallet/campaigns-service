package errors

import (
	"google.golang.org/grpc/codes"
)

type FailedPrecondition struct {
	*richError
}

func NewFailedPrecondition(message string) RichError {
	return &FailedPrecondition{
		richError: newRichError(message),
	}
}

func (f *FailedPrecondition) Error() string {
	return f.message
}

func (f *FailedPrecondition) ToGRPCError() error {
	return f.toGRPCErrorWith(codes.FailedPrecondition)
}
