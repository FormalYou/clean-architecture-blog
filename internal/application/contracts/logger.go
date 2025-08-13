package contracts

// Logger defines the interface for a structured logger.
// This allows for interchangeable logging implementations.
type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	With(args ...interface{}) Logger
}