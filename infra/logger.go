package infra

import (
	"github.com/sirupsen/logrus"
)

func InitLogger(cfg *Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.Server.LogLevel))
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

func getLogLevel(level string) logrus.Level {
	if logLevel, err := logrus.ParseLevel(level); err == nil {
		return logLevel
	}

	return logrus.InfoLevel
}
