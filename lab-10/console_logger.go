package main

import (
	"fmt"
	"time"
)

type ConsoleLogger struct {
	prefix string
}

func NewConsoleLogger(prefix string) *ConsoleLogger {
	return &ConsoleLogger{prefix: prefix}
}

func (l *ConsoleLogger) Info(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Printf("[INFO]  %s [%s] %s\n", timestamp, l.prefix, msg)
}

func (l *ConsoleLogger) Warn(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Printf("[WARN]  %s [%s] %s\n", timestamp, l.prefix, msg)
}

func (l *ConsoleLogger) Error(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Printf("[ERROR]  %s [%s] %s\n", timestamp, l.prefix, msg)
}
