package log

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	AuditLogger *zap.Logger
	// AtomicLevel is a global atomic level for dynamically changing the log level.
	AtomicLevel zap.AtomicLevel
)

func InitLogger() {
	logLevel := viper.GetString("logger.level")
	var initialLevel zapcore.Level
	if err := initialLevel.UnmarshalText([]byte(logLevel)); err != nil {
		initialLevel = zap.InfoLevel
	}
	AtomicLevel = zap.NewAtomicLevelAt(initialLevel)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// File logger
	if viper.IsSet("logger.file.filename") {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   viper.GetString("logger.file.filename"),
			MaxSize:    viper.GetInt("logger.file.max_size"),
			MaxBackups: viper.GetInt("logger.file.max_backups"),
			MaxAge:     viper.GetInt("logger.file.max_age"),
			Compress:   viper.GetBool("logger.file.compress"),
		})

		// Get encoder based on config
		var encoder zapcore.Encoder
		if viper.GetString("logger.encoding") == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, AtomicLevel))
	}

	// Console logger
	consoleWriter := zapcore.AddSync(os.Stdout)
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		consoleWriter,
		AtomicLevel,
	))

	core := zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller())
}

func GetLogger() *zap.Logger {
	return logger
}
