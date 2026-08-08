package main

import (
	"fmt"
	"os"
	"time"
)

type FileLogger struct {
	filename string
	file     *os.File
}

func NewFileLogger(filename string) (*FileLogger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{filename: filename, file: f}, nil
}

func (l *FileLogger) Info(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Fprintf(l.file, "[INFO]  %s  %s\n", timestamp, msg)
}

func (l *FileLogger) Warn(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Fprintf(l.file, "[WARN]  %s  %s\n", timestamp, msg)
}

func (l *FileLogger) Error(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Fprintf(l.file, "[ERROR]  %s  %s\n", timestamp, msg)
}

func (l *FileLogger) Close() {
	l.file.Close()
}
