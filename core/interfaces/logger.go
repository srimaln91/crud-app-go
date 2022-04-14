package interfaces

import "context"

type Logger interface {
	Fatal(ctx context.Context, message string, params ...interface{})
	Error(ctx context.Context, message string, params ...interface{})
	Warn(ctx context.Context, message string, params ...interface{})
	Debug(ctx context.Context, message string, params ...interface{})
	Info(ctx context.Context, message string, params ...interface{})
}
