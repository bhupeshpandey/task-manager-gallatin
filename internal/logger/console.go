package logger

import (
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"log"
)

type consoleLogger struct {
	LogLevel    string
	Environment string
}

func newConsoleLogger(cfg models.LoggingConfig) models.Logger {
	return &consoleLogger{
		LogLevel:    cfg.LogLevel,
		Environment: cfg.Environment,
	}
}

func (l *consoleLogger) Log(logLevel models.LogLevel, message string, logs ...interface{}) {
	log.Println(logs...)
}
