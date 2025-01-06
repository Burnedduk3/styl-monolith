package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

// Logger is a wrapper around the logrus.Logger for centralized logging.
type Logger struct {
	Instance *logrus.Logger
}

// loggerInstance holds the singleton instance of the Logger.
// once ensures that the Logger instance is initialized only once.
var (
	loggerInstance *Logger
	once           sync.Once
)

// GetLogger returns a singleton instance of a configured logrus.Logger.
// The logger's log level is set based on the provided logrus.Level parameter.
// This function ensures the logger is initialized only once using sync.Once.
// Additional custom hooks and text formatting are applied to the logger.
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

// AddEntriesToLogHook is a logrus hook that adds additional data to log entries.
type AddEntriesToLogHook struct{}

// Levels returns all log levels supported by the hook.
func (h *AddEntriesToLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is triggered when a log event occurs and adds the serialized log size to the log entry's data.
func (h *AddEntriesToLogHook) Fire(entry *logrus.Entry) error {
	serialized, err := entry.String()
	if err != nil {
		return err
	}
	size := len(serialized)
	entry.Data["log_size"] = size
	return nil
}
