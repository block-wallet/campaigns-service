package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError(t *testing.T) {
	notFoundError := NewNotFound("not found error")

	assert.NotNil(t, notFoundError)
}

func TestErrorShouldReturnsTheMessage(t *testing.T) {
	notFoundError := NewNotFound("not found error")

	message := notFoundError.Error()

	assert.Equal(t, "not found error", message)
}

func TestToGRPCErrorShouldReturnsNotFoundError(t *testing.T) {
	notFoundError := NewNotFound("not found error")

	err := notFoundError.ToGRPCError()

	assert.EqualError(t, err, "rpc error: code = NotFound desc = not found error")
}
