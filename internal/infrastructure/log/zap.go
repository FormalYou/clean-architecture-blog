package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	logLevel := viper.GetString("logger.level")
	if logLevel != "" {
		level, err := zapcore.ParseLevel(logLevel)
		if err == nil {
			config.Level = zap.NewAtomicLevelAt(level)
		}
	}

	encoding := viper.GetString("logger.encoding")
	if encoding != "" {
		config.Encoding = encoding
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return logger
}