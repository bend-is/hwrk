package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func New(level, format string) (*Logger, error) {
	config := zap.NewProductionConfig()

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	config.Level = lvl
	config.Encoding = format

	log, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{log: log}, nil
}

func (l *Logger) Info(msg string) {
	l.log.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.log.Error(msg)
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}
