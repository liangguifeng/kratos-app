package applog

import (
	"errors"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

// 创建日志对象
func CreateLogger(logConf LoggerConf) (LoggerContextIface, error) {
	if logConf.LogFilePath == "" {
		return nil, errors.New("Log Path is empty")
	}
	appLogger := newLoggerContext()
	appLogger.loggerConf = logConf
	fileWriter, err := newFileWriter(logConf.LogFilePath)
	if err != nil {
		return nil, err
	}

	var ioLogger zerolog.Logger
	env := strings.ToLower(os.Getenv("GO_ENV"))
	if env == "dev" {
		ioLogger = zerolog.New(os.Stderr)
	} else {
		ioLogger = zerolog.New(fileWriter)
	}
	ioLogger.Level(zerolog.Level(appLogger.loggerConf.Level))

	appLogger.logger = &ioLogger
	return appLogger, nil
}
