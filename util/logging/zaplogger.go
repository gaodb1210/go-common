package logging

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultRotateMaxSize    = 100
	defaultRotateMaxBackups = 10
	defaultRotateMaxAge     = 7

	encodeTimeFormat = "2006-01-02 15:04:05.000"
)

// CreateLogger create a zap logger
func CreateLogger(filePath string, compress bool, stats bool, verbose bool) (*zap.Logger, zap.AtomicLevel, error) {
	// 1. init file rotate write
	rotateWriter := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    defaultRotateMaxSize,
		MaxAge:     defaultRotateMaxAge,
		MaxBackups: defaultRotateMaxBackups,
		LocalTime:  true,
		Compress:   compress,
	}
	syncer := zapcore.AddSync(rotateWriter)
	// 2. init encoder, set time format, set output type, json or console
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(encodeTimeFormat)
	encoderConfig.ConsoleSeparator = " | "

	// 3. init log level
	var level = zap.NewAtomicLevel()
	if verbose {
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	// 4. init core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		syncer,
		level,
	)
    // set log options, caller, stack and so on
	var opts []zap.Option
	if !stats {
		opts = append(opts, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel), zap.AddCallerSkip(1))
	}
	// 5. create logger
	return zap.New(core, opts...), level, nil
}
