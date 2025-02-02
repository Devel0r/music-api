package logger

import (
    "github.com/sirupsen/logrus"
)

type Logger struct {
    logger *logrus.Logger
}

func NewLogger() *Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.TextFormatter{})

    return &Logger{
        logger: logger,
    }
}

func (l *Logger) Debug(msg string, fields logrus.Fields) {
    l.logger.WithFields(fields).Debug(msg)
}

func (l *Logger) Info(msg string, fields logrus.Fields) {
    l.logger.WithFields(fields).Info(msg)
}

func (l *Logger) Error(msg string, fields logrus.Fields) {
    l.logger.WithFields(fields).Error(msg)
}