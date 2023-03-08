package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RichError interface {
	error
	ToGRPCError() error
}

type richError struct {
	message string
}

func (r *richError) toGRPCErrorWith(code codes.Code) error {
	return status.Error(code, r.message)
}

func newRichError(message string) *richError {
	return &richError{message: message}
}
