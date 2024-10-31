package logging

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func GetInstance() *slog.Logger {
	return logger
}

func init() {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}
}
