package logger

import (
	"context"

	"go.uber.org/zap"
)

type ContextKey string

const (
	messageIDField ContextKey = "message_id"
)

type sugaredLogger struct {
	*zap.SugaredLogger
	messageIDField ContextKey
}

func NewSugaredLogger(zapSugaredLogger *zap.SugaredLogger, messageIDField ContextKey) *sugaredLogger {
	return &sugaredLogger{
		zapSugaredLogger,
		messageIDField,
	}
}

func (s *sugaredLogger) GetMessageIDField() ContextKey {
	return s.messageIDField
}

func (s *sugaredLogger) WithCtx(ctx context.Context, args ...interface{}) *zap.SugaredLogger {
	// check if ctx is not nil to avoid panics.
	if ctx != nil {
		if requestID := ctx.Value(s.messageIDField); requestID != nil {
			args = append(args, string(s.messageIDField), requestID)
		}
	}

	// args can't be nil or empty on with, if they are nil, it returns the zap sugared logger.
	if len(args) == 0 {
		return s.SugaredLogger
	}

	return s.With(args...)
}
