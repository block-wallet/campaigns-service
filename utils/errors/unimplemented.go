package errors

import (
	"google.golang.org/grpc/codes"
)

type unimplemented struct {
	*richError
}

func NewUnimplemented(message string) RichError {
	return &unimplemented{
		newRichError(message),
	}
}

func (u *unimplemented) Error() string {
	return u.message
}

func (u *unimplemented) ToGRPCError() error {
	return u.toGRPCErrorWith(codes.Unimplemented)
}
