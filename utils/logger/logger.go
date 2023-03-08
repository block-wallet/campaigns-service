package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Sugar logger to use on the project
var Sugar *sugaredLogger

// init is used to avoid panics on testing.
// for production environments you have to call the Initialize function .
func init() {
	config, _ := zap.NewProductionConfig().Build()
	Sugar = NewSugaredLogger(config.Sugar(), messageIDField)
}

func Initialize(level, version string) error {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(getZapLevel(level))
	config.InitialFields = map[string]interface{}{
		"version": version,
	}
	logger, err := config.Build()
	if err != nil {
		return err
	}

	Sugar = NewSugaredLogger(logger.Sugar(), messageIDField)
	return nil
}

func getZapLevel(level string) zapcore.Level {
	levelMap := map[string]zapcore.Level{
		"debug":   zapcore.DebugLevel,
		"info":    zapcore.InfoLevel,
		"error":   zapcore.ErrorLevel,
		"warning": zapcore.WarnLevel,
		"fatal":   zapcore.FatalLevel,
		"panic":   zapcore.PanicLevel,
	}
	return levelMap[level]
}
