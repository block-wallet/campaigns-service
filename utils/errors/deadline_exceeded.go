package errors

import (
	"google.golang.org/grpc/codes"
)

type deadlineExceeded struct {
	*richError
}

func NewDeadlineExceeded(message string) RichError {
	return &deadlineExceeded{
		newRichError(message),
	}
}

func (d *deadlineExceeded) Error() string {
	return d.message
}

func (d *deadlineExceeded) ToGRPCError() error {
	return d.toGRPCErrorWith(codes.DeadlineExceeded)
}
