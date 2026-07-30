package main

import "fmt"

// Events in the system
type LoginEvent struct {
	UserID    int
	IPAddress string
	Success   bool
}

type PurchaseEvent struct {
	UserID    int
	ProductID string
	Amount    float64
}

type ErrorEvent struct {
	Code    int
	Message string
	Fatal   bool
}

// Process handles any event type
func Process(event any) string {
	// Use a type switch to handle each type differently
	switch v := event.(type) {
	// LoginEvent:   "LOGIN user=42 ip=192.168.1.1 success=true"
	case LoginEvent:
		return fmt.Sprintf("LOGIN user=%d ip=%s success=%t", v.UserID, v.IPAddress, v.Success)
	// PurchaseEvent: "PURCHASE user=42 product=SKU-001 amount=PKR 1500.00"
	case PurchaseEvent:
		return fmt.Sprintf("PURCHASE user=%d product=%s amount=PKR %.2f", v.UserID, v.ProductID, v.Amount)
	// ErrorEvent:   "ERROR [404] page not found (fatal=false)"
	case ErrorEvent:
		return fmt.Sprintf("ERROR [%d] %s (fatal=%t)", v.Code, v.Message, v.Fatal)
		// unknown:      "UNKNOWN event type: <type>"
	default:
		return fmt.Sprintf("UNKNOWN event type: <%T>", v)
	}
}

// ProcessBatch processes a slice of mixed events
func ProcessBatch(events []any) {
	for i, e := range events {
		fmt.Printf("[%d] %s\n", i+1, Process(e))
	}
}

func main() {
	events := []any{
		LoginEvent{UserID: 1, IPAddress: "192.168.1.1", Success: true},
		PurchaseEvent{UserID: 1, ProductID: "SKU-001", Amount: 1500},
		LoginEvent{UserID: 2, IPAddress: "10.0.0.5", Success: false},
		ErrorEvent{Code: 500, Message: "database timeout", Fatal: true},
		"unexpected string", // unknown type
	}
	ProcessBatch(events)
}
