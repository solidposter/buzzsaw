package main

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func slogsetup(logfile string, debug bool) {

	var logger *slog.Logger
	logWriter := &lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    100, // megabytes
		MaxBackups: 5,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	if debug {
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		//handler := slog.NewJSONHandler(os.Stdout, opts)
		handler := slog.NewJSONHandler(logWriter, opts)
		logger = slog.New(handler)
	} else {
		logger = slog.New(slog.NewJSONHandler(logWriter, nil))
	}

	slog.SetDefault(logger)
}
