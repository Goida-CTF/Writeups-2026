package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level zapcore.Level) (*zap.Logger, error) {
	lc := zap.NewProductionConfig()
	lc.Level = zap.NewAtomicLevelAt(level)
	lc.OutputPaths = []string{"stdout"}

	l, err := lc.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("lc.Build: %w", err)
	}
	return l, nil
}
