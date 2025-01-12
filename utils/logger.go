package utils

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}