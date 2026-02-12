package infra

import (
	"sync"

	"go.uber.org/zap"
)

var (
	loggerOnce sync.Once
	loggerInst *zap.Logger
	loggerErr  error
)

func GetLogger() (*zap.Logger, error) {
	loggerOnce.Do(func() {
		loggerInst, loggerErr = zap.NewProduction()
	})
	return loggerInst, loggerErr
}
