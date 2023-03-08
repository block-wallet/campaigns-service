package errors

import (
	"google.golang.org/grpc/codes"
)

type permissionDenied struct {
	*richError
}

func NewPermissionDenied(message string) RichError {
	return &permissionDenied{
		newRichError(message),
	}
}

func (p *permissionDenied) Error() string {
	return p.message
}

func (p *permissionDenied) ToGRPCError() error {
	return p.toGRPCErrorWith(codes.PermissionDenied)
}
