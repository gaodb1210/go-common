package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

const (
	encodeTimeFormat        = "2006-01-02 15:04:05.000"
	defaultRotateMaxSize    = 1024
	defaultRotateMaxDays    = 1
	defaultRotateMaxBackups = 14
)

func createLogger(fileName, level string, rotateDays, saveDays int, compress bool) (*zap.Logger, error) {
	if rotateDays == 0 {
		rotateDays = defaultRotateMaxDays
	}
	if saveDays == 0 {
		saveDays = defaultRotateMaxBackups
	}
	// init file rotate
	rotateWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    defaultRotateMaxSize,
		MaxAge:     rotateDays,
		MaxBackups: saveDays,
		LocalTime:  true,
		Compress:   compress,
	}

	// init encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(encodeTimeFormat)
	encoderConfig.ConsoleSeparator = " | "

	// init level
	var l = zap.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		l = zap.DebugLevel
	case "warn", "warning":
		l = zap.WarnLevel
	case "err", "error":
		l = zap.ErrorLevel
	}

	// init core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(rotateWriter),
		zap.NewAtomicLevelAt(l),
	)

	// create log
	return zap.New(core,
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(1)), nil
}
