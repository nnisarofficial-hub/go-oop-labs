package main

type NoopLogger struct{}

func (l *NoopLogger) Info(msg string)  {}
func (l *NoopLogger) Warn(msg string)  {}
func (l *NoopLogger) Error(msg string) {}
