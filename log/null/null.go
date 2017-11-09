package null

import (
	"github.com/gliderlabs/com/objects"
	"github.com/gliderlabs/stdcom/log"
)

type Logger struct{}

// Register the null logger with a registry
func Register(registry *objects.Registry) error {
	return registry.Register(&objects.Object{
		Value: &Logger{},
	})
}

func (l *Logger) With(args ...interface{}) log.Logger {
	return l
}
func (l *Logger) Debug(args ...interface{})                       {}
func (l *Logger) Debugf(template string, args ...interface{})     {}
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {}
func (l *Logger) Info(args ...interface{})                        {}
func (l *Logger) Infof(template string, args ...interface{})      {}
func (l *Logger) Infow(msg string, keysAndValues ...interface{})  {}
func (l *Logger) Error(args ...interface{})                       {}
func (l *Logger) Errorf(template string, args ...interface{})     {}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {}
