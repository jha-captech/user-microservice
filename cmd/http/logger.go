package main

import (
	"log/slog"
	"os"
)

func newLogger(useJSON bool) *slog.Logger {
	logger := slog.Default()
	if useJSON {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return logger
}
