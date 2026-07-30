package main

import (
	"fmt"
	"os"
	"time"
)

const (
	format = "2006/01/02 15:04:05"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

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
	fmt.Fprintf(l.file, "[INFO]  %s [orders] %s\n", timestamp, msg)
}
func (l *FileLogger) Warn(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Printf("[WARN]  %s [orders] %s\n", timestamp, msg)
}
func (l *FileLogger) Error(msg string) {
	timestamp := time.Now().Format(format)
	fmt.Printf("[ERROR]  %s [orders] %s\n", timestamp, msg)
}
func (l *FileLogger) Close() {
	l.file.Close()
}

type NoopLogger struct{}

func (l *NoopLogger) Info(msg string)  {}
func (l *NoopLogger) Warn(msg string)  {}
func (l *NoopLogger) Error(msg string) {}

type OrderService struct {
	logger Logger
	orders map[int]string
}

func NewOrderService(logger Logger) *OrderService {
	return &OrderService{logger: logger, orders: map[int]string{}}
}

func (s *OrderService) PlaceOrder(id int, item string) {
	s.logger.Info(fmt.Sprintf("placing order %d for item: %s", id, item))
	s.orders[id] = item
	s.logger.Info(fmt.Sprintf("order %d placed successfully", id))
}

func (s *OrderService) CancelOrder(id int) error {
	if _, exists := s.orders[id]; !exists {
		s.logger.Warn(fmt.Sprintf("attempted to cancel non-existent order %d", id))
		return fmt.Errorf("order %d not found", id)
	}
	delete(s.orders, id)
	s.logger.Info(fmt.Sprintf("order %d cancelled", id))
	return nil
}

func main() {
	svc := NewOrderService(NewConsoleLogger("orders"))
	svc.PlaceOrder(1, "Barbari Buck — Premium Grade")
	svc.CancelOrder(999)
	testSvc := NewOrderService(&NoopLogger{})
	testSvc.PlaceOrder(2, "test item")
}
