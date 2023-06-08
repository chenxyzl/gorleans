package glog

import (
	"fmt"
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
	return sugar.Desugar().WithOptions(zap.AddCallerSkip(-1)).WithOptions(opts...).Sugar()
}

func With(args ...interface{}) {
	sugar = sugar.With(args...)
}

func NewWith(args ...interface{}) *zap.SugaredLogger {
	return sugar.Desugar().WithOptions(zap.AddCallerSkip(-1)).Sugar().With(args...)
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
func Errorf(template string, args ...interface{}) error {
	err := fmt.Errorf(template, args...)
	sugar.Error(err)
	return err
}

// Fatalf uses fmt.Sprintf to log a templated message, then panics.
func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
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
