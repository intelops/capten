package clog

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	Logger = GetLogger()
)

type cliFormatter struct {
}

func (f *cliFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor *color.Color
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = color.New(color.FgGreen)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgRed, color.Bold)
	case logrus.DebugLevel:
		levelColor = color.New(color.FgCyan)
	default:
		levelColor = color.New()
	}
	message := fmt.Sprintf("[%s] %s\n", levelColor.Sprint(strings.ToUpper(entry.Level.String())), entry.Message)
	return []byte(message), nil
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	level := os.Getenv("LOG_LEVEL")
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)
	logger.SetFormatter(&cliFormatter{})
	return logger
}
