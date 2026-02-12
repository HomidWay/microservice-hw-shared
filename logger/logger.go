package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log Logger

	once sync.Once
)

func GetLogger() Logger {
	once.Do(func() {
		newLogger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		log = newLogger.Sugar()
	})

	return log
}
