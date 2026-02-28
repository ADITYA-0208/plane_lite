package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a production-style zap logger. Call Sync() before exit.
func New(development bool) (*zap.Logger, error) {
	var cfg zap.Config
	if development {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	return cfg.Build()
}
