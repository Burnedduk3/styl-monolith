package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Logger struct {
	instance *logrus.Logger
}

var (
	loggerInstance *Logger
	once           sync.Once
)

func GetLogger(level logrus.Level) *logrus.Logger {
	once.Do(func() {
		log := logrus.New()
		log.SetReportCaller(true)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		log.SetLevel(level)

		loggerInstance = &Logger{instance: log}
	})

	return loggerInstance.instance
}
