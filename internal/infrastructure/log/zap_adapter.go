package log

import (
	"github.com/FormalYou/clean-architecture-blog/internal/application/contracts"
	"go.uber.org/zap"
)

// ZapAdapter is an adapter for the zap logger to conform to the service.Logger interface.
type ZapAdapter struct {
	logger *zap.SugaredLogger
}

// NewZapAdapter creates a new ZapAdapter.
func NewZapAdapter(logger *zap.Logger) contracts.Logger {
	return &ZapAdapter{logger: logger.Sugar()}
}

// Info logs an info message.
func (a *ZapAdapter) Info(msg string, args ...interface{}) {
	a.logger.Infow(msg, args...)
}

// Warn logs a warning message.
func (a *ZapAdapter) Warn(msg string, args ...interface{}) {
	a.logger.Warnw(msg, args...)
}

// Error logs an error message.
func (a *ZapAdapter) Error(msg string, args ...interface{}) {
	a.logger.Errorw(msg, args...)
}

// With adds structured context to a logger.
func (a *ZapAdapter) With(args ...interface{}) contracts.Logger {
	return &ZapAdapter{logger: a.logger.With(args...)}
}
