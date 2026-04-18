package logger

import (
	"fmt"
	"time"
)

// Farben
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func logWithPrefix(color, level, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(color + "[" + level + "] " + timestamp + " - " + msg + ColorReset)
}

func Info(msg string) {
	logWithPrefix(ColorBlue, "INFO", msg)
}

func Success(msg string) {
	logWithPrefix(ColorGreen, "SUCCESS", msg)
}

func Warning(msg string) {
	logWithPrefix(ColorYellow, "WARNING", msg)
}

func Error(msg string) {
	logWithPrefix(ColorRed, "ERROR", msg)
}