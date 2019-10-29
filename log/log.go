package log

import (
	"fmt"
	"log" // nolint
	"os"

	"github.com/breathman/go-dig-services/config"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Service struct {
	*logrus.Logger
}

type CtxLogger struct {
	*logrus.Entry
}

func NewLogger(config *config.APPConfig) *Service {
	logConfig := config.Main.Logger
	logger := Service{}
	logger.Logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &prefixed.TextFormatter{
		ForceColors:     logConfig.IsDevMode,
		ForceFormatting: logConfig.IsDevMode,
	}
	switch logConfig.LoggerLevel {
	case "error":
		logger.Level = logrus.ErrorLevel
	case "info":
		logger.Level = logrus.InfoLevel
	default:
		logger.Level = logrus.DebugLevel
	}
	return &logger
}

func (ls *Service) NewPrefix(prefix string) *CtxLogger {
	ctxLog := CtxLogger{
		Entry: ls.WithField("prefix", prefix),
	}
	return &ctxLog
}

func (ctxl *CtxLogger) NewPrefix(prefix string) *CtxLogger {
	ctxLog := CtxLogger{
		Entry: ctxl.WithField("prefix", prefix),
	}
	return &ctxLog
}

func (ctxl *CtxLogger) Print(v ...interface{}) {
	ctxl.Debug(v...)
}

func (ctxl *CtxLogger) AddPrefix(prefix string) *CtxLogger {
	var newPrefix string
	if data, ok := ctxl.Data["prefix"]; !ok {
		newPrefix = prefix
	} else {
		newPrefix = fmt.Sprintf("%s.%s", data, prefix)
	}
	return &CtxLogger{
		Entry: ctxl.WithField("prefix", newPrefix),
	}
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
