package util

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Logger is an adapter type for zap's SugaredLogger
type Logger struct {
	logger *zap.SugaredLogger
}

// NewLogger creates a new Logger.
func NewLogger(level, serviceName string) (*Logger, error) {
	atom := zap.NewAtomicLevel()

	switch level {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "warn":
		atom.SetLevel(zap.WarnLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	default:
		level = "info"
		atom.SetLevel(zap.InfoLevel)
	}

	cfg := zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		Encoding:          "json",
		ErrorOutputPaths:  []string{"stdout"},
		Level:             atom,
		OutputPaths:       []string{"stdout"},
	}

	if cfg.InitialFields == nil {
		cfg.InitialFields = make(map[string]interface{})
	}

	cfg.InitialFields["service"] = serviceName

	l, err := cfg.Build()

	if err != nil {
		return nil, errors.Wrap(err, "Unable to initialize zap logger")
	}

	l.Debug("Logger Created", zap.String("level", level))

	return &Logger{logger: l.Sugar()}, nil
}

// LoggerMiddleware is a decorator for a HTTP Request, adding structured logging functionality
func LoggerMiddleware(inner http.Handler, logger *Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := RequestIDLogger(logger, r)
		log.Debugw("request received",
			"address", r.RemoteAddr,
			"method", r.Method,
			"path", r.RequestURI,
		)
		inner.ServeHTTP(w, r)
		log.Infow("request completed",
			"address", r.RemoteAddr,
			"method", r.Method,
			"path", r.RequestURI,
			"duration", time.Since(start),
		)
	})
}

// RequestIDLogger extracts the requestID from the request context, adds it to
// a child logger and then returns the child logger.
func RequestIDLogger(l *Logger, r *http.Request) *Logger {
	reqID := RequestIDFromContext(r.Context())
	log := l.logger.With("requestID", reqID)
	return &Logger{log}
}

// RequestIDLoggerFromContext extract the requestID from a passed context, adds
// it to a child logger and then returns the child logger.
func RequestIDLoggerFromContext(ctx context.Context, l *Logger) *Logger {
	reqID := RequestIDFromContext(ctx)
	log := l.logger.With("requestID", reqID)
	return &Logger{log}

}

// Info uses fmt.Sprint to log a templated message.
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof uses fmt.Sprint to log a templated message.
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

// Infow uses fmt.Sprint to log a templated message.
func (l *Logger) Infow(msg string, kv ...interface{}) {
	l.logger.Infow(msg, kv...)
}

// Error uses fmt.Sprint to log a templated message.
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Errorw(msg string, kv ...interface{}) {
	l.logger.Errorw(msg, kv...)
}

// Warn uses fmt.Sprint to log a templated message.
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

// Warnw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (l *Logger) Warnw(msg string, kv ...interface{}) {
	l.logger.Warnw(msg, kv...)
}

// Debug uses fmt.Sprint to log a templated message.
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

// Debugw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (l *Logger) Debugw(msg string, kv ...interface{}) {
	l.logger.Debugw(msg, kv...)
}

// Panic uses fmt.Sprint to log a templated message, then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(msg string, args ...interface{}) {
	l.logger.Panicf(msg, args...)
}

// Panicw logs a message with some additional context, then panics The variadic key-value pairs are treated as they are in With.
func (l *Logger) Panicw(msg string, kv ...interface{}) {
	l.logger.Panicw(msg, kv...)
}

// Fatal uses fmt.Sprint to log a templated message, then calls os.Exit.
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The variadic key-value pairs are treated as they are in With.
func (l *Logger) Fatalw(msg string, kv ...interface{}) {
	l.logger.Fatalw(msg, kv...)
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() {
	l.logger.Sync()
}
