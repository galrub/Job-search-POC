package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LOG zerolog.Logger

func InitStaticLogger() {
	if os.Getenv("LOG_TO_CONSOLE") == "true" {
		fmt.Print("Using Std. out for logging")
		LOG = zerolog.New(os.Stdout)
		return
	}
	logFile := os.Getenv("LOG_FILE")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	logger := zerolog.New(lumberjackLogger).With().Timestamp().Logger()
	LOG = logger
}
