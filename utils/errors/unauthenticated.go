package errors

import (
	"google.golang.org/grpc/codes"
)

type unauthenticated struct {
	*richError
}

func NewUnauthenticated(message string) RichError {
	return &unauthenticated{
		newRichError(message),
	}
}

func (u *unauthenticated) Error() string {
	return u.message
}

func (u *unauthenticated) ToGRPCError() error {
	return u.toGRPCErrorWith(codes.Unauthenticated)
}
