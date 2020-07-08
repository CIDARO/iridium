package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// SetupLogger creates the default logger
func SetupLogger() {
	Logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{DisableColors: false, FullTimestamp: true},
		Level:     logrus.InfoLevel,
	}
}
