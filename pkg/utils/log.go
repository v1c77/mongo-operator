package utils

import (
	"fmt"
	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type LogHelper struct {
	l logr.Logger
}

func NewLogger(name string) *LogHelper {
	logger := logf.Log.WithName(name)

	return &LogHelper{
		l: logger,
	}
}

func (l *LogHelper) Infof(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
	l.l.Info(msg)
}

func (l *LogHelper) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
	l.l.Info(msg)
}

func (l *LogHelper) L() logr.Logger {
	return l.l
}

func (l *LogHelper) Info(msg string, keysAndValues ...interface{}) {
	l.l.Info(msg, keysAndValues...)
}

func (l *LogHelper) Error(err error, msg string, keysAndValues ...interface{}) {
	l.l.Error(err, msg, keysAndValues...)
}
