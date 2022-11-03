package cmd

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func getLogger() logr.Logger {
	// setup logger
	zapConf := zap.NewDevelopmentConfig()
	zapConf.DisableStacktrace = true // due to output wrapped error in errorVerbose
	logger, _ := zapConf.Build()
	return zapr.NewLogger(logger)
}
