package main

import (
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func slogsetup(logfile string, debug bool) {
	var logger *slog.Logger
	var opts *slog.HandlerOptions

	if debug {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	} else {
		opts = nil
	}

	if logfile == "" {
		handler := slog.NewJSONHandler(os.Stdout, opts)
		logger = slog.New(handler)
	} else {
		logWriter := &lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    100, // megabytes
			MaxBackups: 5,
			MaxAge:     28,    //days
			Compress:   false, // disabled by default
		}
		handler := slog.NewJSONHandler(logWriter, opts)
		logger = slog.New(handler)
	}

	slog.SetDefault(logger)
}
