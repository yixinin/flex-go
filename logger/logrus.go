package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	log *logrus.Logger
}

func FromLogrus(log *logrus.Logger) *LogrusLogger {
	return &LogrusLogger{
		log: log,
	}
}

func (logger *LogrusLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Debugf(msg, args...)
}
func (logger *LogrusLogger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Warnf(msg, args...)
}
func (logger *LogrusLogger) Infof(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Infof(msg, args...)
}
func (logger *LogrusLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Errorf(msg, args...)
}

func (logger *LogrusLogger) Debug(ctx context.Context, args ...interface{}) {
	logger.log.Debug(args...)
}
func (logger *LogrusLogger) Warn(ctx context.Context, args ...interface{}) {
	logger.log.Warn(args...)
}
func (logger *LogrusLogger) Info(ctx context.Context, args ...interface{}) {
	logger.log.Info(args...)
}
func (logger *LogrusLogger) Error(ctx context.Context, args ...interface{}) {
	logger.log.Error(args...)
}
