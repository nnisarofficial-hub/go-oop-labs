package main

import "fmt"

type ConsoleLogger struct{}

func (c *ConsoleLogger) Info(msg string)  { fmt.Println("[INFO]", msg) }
func (c *ConsoleLogger) Warn(msg string)  { fmt.Println("[WARN]", msg) }
func (c *ConsoleLogger) Error(msg string) { fmt.Println("[ERROR]", msg) }
