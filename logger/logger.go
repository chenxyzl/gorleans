package logger

import (
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
	sugar = l.Sugar()
}

func WithOptions(opts ...zap.Option) {
	// sugar= sugar.Desugar().WithOptions(zap.IncreaseLevel(level)).Sugar()
	sugar = sugar.Desugar().WithOptions(opts...).Sugar()
}

func NewWithOptions(opts ...zap.Option) *zap.SugaredLogger {
	return sugar.Desugar().WithOptions(opts...).Sugar()
}

func WithFields(field ...zap.Field) {
	sugar = sugar.Desugar().With(field...).Sugar()
}

func NewWithFields(field ...zap.Field) *zap.SugaredLogger {
	return sugar.Desugar().With(field...).Sugar()
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Error uses fmt.Sprintf to log a templated message.
func Error(err error) error {
	sugar.Error(err)
	return err
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Panic uses fmt.Sprintf to log a templated message, then panics.
func Panic(err error) error {
	sugar.Panic(err)
	return err
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return sugar.Sync()
}
