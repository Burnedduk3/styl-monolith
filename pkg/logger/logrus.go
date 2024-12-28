package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Logger struct {
	Instance *logrus.Logger
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
			DisableColors:          false,
			FullTimestamp:          true,
			DisableLevelTruncation: false,
		})
		log.AddHook(&AddEntriesToLogHook{})
		log.SetLevel(level)
		loggerInstance = &Logger{Instance: log}
	})

	return loggerInstance.Instance
}

type AddEntriesToLogHook struct{}

func (h *AddEntriesToLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *AddEntriesToLogHook) Fire(entry *logrus.Entry) error {
	serialized, err := entry.String()
	if err != nil {
		return err
	}
	size := len(serialized)
	entry.Data["log_size"] = size
	return nil
}
