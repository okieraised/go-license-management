package logging

import (
	"fmt"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	logger *zap.Logger
}

var currentLogger = &Logger{}
var childLogger *zap.Logger
var logConfig zap.Config
var err error

func init() {
	logConfig = zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	logConfig.EncoderConfig = ecszap.ECSCompatibleEncoderConfig(logConfig.EncoderConfig)
	childLogger, err = logConfig.Build()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		os.Exit(1)
	}
}

func NewDefaultLogger() {
	newCLogger, err := logConfig.Build()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		os.Exit(1)
	}
	currentLogger.logger = newCLogger
}

func NewECSLogger() *Logger {
	cLogger, err := logConfig.Build()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		os.Exit(1)
	}

	return &Logger{logger: cLogger}
}

func GetInstance() *Logger {
	return currentLogger
}

func (l *Logger) GetLogger() *zap.Logger {
	return l.logger
}

func (l *Logger) GetSugarLogger() *zap.SugaredLogger {
	return l.logger.Sugar()
}

func (l *Logger) WithCustomFields(fields ...zap.Field) *zap.Logger {
	l.logger = childLogger
	l.logger = l.logger.With(fields...)
	return l.logger
}

func (l *Logger) WithCustomStringFields(k string, v string) *zap.Logger {
	l.logger = childLogger
	l.logger = l.logger.With(zap.String(k, v))
	return l.logger
}
