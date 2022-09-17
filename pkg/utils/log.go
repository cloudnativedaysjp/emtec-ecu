package utils

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func GetLogger(ctx context.Context) logr.Logger {
	logger, err := logr.FromContext(ctx)
	if err != nil {
		zaplogger, _ := zap.NewDevelopment()
		logger = zapr.NewLogger(zaplogger)
	}
	return logger
}
