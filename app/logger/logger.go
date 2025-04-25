package logger

import (
	"fmt"
	"time"
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func Success(msg string) {
	logWithColor("SUCCESS", msg, green)
}

func Info(msg string) {
	logWithColor("INFO", msg, yellow)
}

func Error(msg string) {
	logWithColor("ERROR", msg, red)
}

func logWithColor(prefix, msg, color string) {
	fmt.Printf("[%s] %s%s%s: %s\n",
		time.Now().Format("15:04:05"),
		color, prefix, reset,
		msg,
	)
}
