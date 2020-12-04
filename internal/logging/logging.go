package logging

import (
	"github.com/sirupsen/logrus"
)

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

type registeredLogger struct {
	logger       StandardLogger
	isRegistered bool
}

var (
	globalLogger = registeredLogger{StandardLogger{}, false}
)

func SetGlobalLogger(logger StandardLogger) {
	globalLogger = registeredLogger{logger, true}
}

func GlobalLogger() StandardLogger {
	return globalLogger.logger
}

var standardFields logrus.Fields

// NewLogger initializes the standard logger
func NewLogger(application string) *StandardLogger {
	standardFields = logrus.Fields{
		"app": application,
	}

	var baseLogger = logrus.New()

	var standardLogger = &StandardLogger{baseLogger}

	standardLogger.Formatter = &logrus.JSONFormatter{}
	standardLogger.SetLevel(logrus.DebugLevel)
	return standardLogger
}

func (l *StandardLogger) Debug(fields logrus.Fields, msg string) {
	l.WithFields(standardFields).WithFields(fields).Debug(msg)
}

func (l *StandardLogger) Info(msg string) {
	l.WithFields(standardFields).Info(msg)
}

func (l *StandardLogger) Warn(msg string) {
	l.WithFields(standardFields).Warn(msg)
}

func (l *StandardLogger) Error(msg string) {
	l.WithFields(standardFields).Error(msg)
}

func (l *StandardLogger) ErrorWErr(msg string, err error) {
	l.WithFields(standardFields).WithError(err).Error(msg)
}
