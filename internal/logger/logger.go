package logger

import "github.com/bhupeshpandey/task-manager-gallatin/internal/models"

func NewLogger(config models.LoggingConfig) models.Logger {
	var loggerInst models.Logger
	switch config.Type {
	case "console":
		loggerInst = newConsoleLogger(config)
	case "zap":
		loggerInst = newZapLogger(config)
	}
	return loggerInst
}
