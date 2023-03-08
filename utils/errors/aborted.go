package errors

import (
	"google.golang.org/grpc/codes"
)

type aborted struct {
	*richError
}

func NewAborted(message string) RichError {
	return &aborted{
		newRichError(message),
	}
}

func (a *aborted) Error() string {
	return a.message
}

func (a *aborted) ToGRPCError() error {
	return a.toGRPCErrorWith(codes.Aborted)
}
