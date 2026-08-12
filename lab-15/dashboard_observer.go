package main

import "fmt"

type DashboardObserver struct{}

func (d *DashboardObserver) Name() string { return "Dashboard" }
func (d *DashboardObserver) OnEvent(event Event) {
	details := ""
	for key, value := range event.Payload {
		details += fmt.Sprintf("%s=%s ", key, value)
	}
	fmt.Printf("  [%s] %s: %s\n", d.Name(), event.Type, details)
}
