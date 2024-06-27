package main

import (
	"log/slog"
	"os"
)

func newLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
