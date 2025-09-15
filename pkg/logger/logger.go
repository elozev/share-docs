package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

type LogConfig struct {
	Level       string
	Environmnet string
	OutputPath  string
	ServiceName string
	Version     string
}

func NewLogger(config LogConfig) (*Logger, error) {
	if config.Level == "" {
		config.Level = "info"
	}
	if config.Environmnet == "" {
		config.Environmnet = "development"
	}
	if config.OutputPath == "" {
		config.OutputPath = "stdout"
	}

	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	var encoderConfig zapcore.EncoderConfig
	if config.Environmnet == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.StacktraceKey = "stacktrace"

	var encoder zapcore.Encoder
	if config.Environmnet == "production" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	if config.OutputPath == "stdout" {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(config.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	zapLogger = zapLogger.With(
		zap.String("service", config.ServiceName),
		zap.String("version", config.Version),
	)

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	return &Logger{
		Logger: l.Logger.With(zapFields...),
		sugar:  l.Logger.With(zapFields...).Sugar(),
	}
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.Any(key, value)),
		sugar:  l.Logger.With(zap.Any(key, value)).Sugar(),
	}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.Error(err)),
		sugar:  l.Logger.With(zap.Error(err)).Sugar(),
	}
}

func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithField("user_id", userID)
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithField("request_id", requestID)
}

// Sugar returns the sugared logger for printf-style logging
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}
