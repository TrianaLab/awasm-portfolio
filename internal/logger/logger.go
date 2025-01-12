package logger

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitLogger(level logrus.Level) {
	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.Level = level
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		InitLogger(logrus.InfoLevel) // Default level
	}
	return logger
}

func Trace(fields logrus.Fields, msg string) {
	GetLogger().WithFields(fields).Trace(msg)
}

func Warn(fields logrus.Fields, msg string) {
	GetLogger().WithFields(fields).Warn(msg)
}

func Info(fields logrus.Fields, msg string) {
	GetLogger().WithFields(fields).Info(msg)
}

func Error(fields logrus.Fields, msg string) {
	GetLogger().WithFields(fields).Error(msg)
}
