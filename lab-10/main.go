package main

import (
	"fmt"
)

const (
	format = "2006/01/02 15:04:05"
)

func main() {
	consoleLogger := NewConsoleLogger("orders")
	_ = consoleLogger
	fileLogger, _ := NewFileLogger("service.log")
	_ = fileLogger
	noOpLogger := NoopLogger{}
	_ = noOpLogger
	svc := NewOrderService(consoleLogger)
	svc.PlaceOrder(1, "Barbari Buck — Premium Grade")
	svc.CancelOrder(999)
	fileLogger, err := NewFileLogger("orders.log")
	if err != nil {
		fmt.Println("failed to create file logger:", err)
		return
	}
	defer fileLogger.Close()
	fileSvc := NewOrderService(fileLogger)
	fileSvc.PlaceOrder(3, "file logged item")

	testSvc := NewOrderService(&NoopLogger{})
	testSvc.PlaceOrder(2, "test item")
}
