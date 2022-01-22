package logger

import "context"

type Interface interface {
	Debugf(ctx context.Context, msg string, args ...interface{})
	Warnf(ctx context.Context, msg string, args ...interface{})
	Infof(ctx context.Context, msg string, args ...interface{})
	Errorf(ctx context.Context, msg string, args ...interface{})

	Debug(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
}

var Logger Interface

func Debugf(ctx context.Context, msg string, args ...interface{}) {
	Logger.Debugf(ctx, msg, args...)
}
func Warnf(ctx context.Context, msg string, args ...interface{}) {
	Logger.Warnf(ctx, msg, args...)
}
func Infof(ctx context.Context, msg string, args ...interface{}) {
	Logger.Infof(ctx, msg, args...)
}
func Errorf(ctx context.Context, msg string, args ...interface{}) {
	Logger.Errorf(ctx, msg, args...)
}
func Debug(ctx context.Context, args ...interface{}) {
	Logger.Debug(ctx, args...)
}
func Warn(ctx context.Context, args ...interface{}) {
	Logger.Warn(ctx, args...)
}
func Info(ctx context.Context, args ...interface{}) {
	Logger.Info(ctx, args...)
}
func Error(ctx context.Context, args ...interface{}) {
	Logger.Error(ctx, args...)
}

func Init(l Interface) {
	Logger = l
}
