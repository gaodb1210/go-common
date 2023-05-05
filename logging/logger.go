package logging

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func InitLogger(filePath, level string, rotateDays, saveDays int, compress bool) error {
	l, err := createLogger(filePath, level, rotateDays, saveDays, compress)
	if err != nil || l == nil {
		return err
	}

	logger = l.Sugar()
	return nil
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args)
}
