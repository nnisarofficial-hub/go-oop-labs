package main

import "fmt"

type Observer interface {
	OnEvent(event Event)
	Name() string
}

func main() {
	dashboard := &DashboardObserver{}
	tracker := &HealthTrackerObserver{}
	alert := &AlertObserver{threshold: "critical"}
	monitor := &AnimalMonitor{AnimalID: "D-001"}
	fmt.Println("Subscribing: Dashboard, HealthTracker, AlertSystem")
	monitor.Subscribe(dashboard)
	monitor.Subscribe(tracker)
	monitor.Subscribe(alert)
	fmt.Println("\nEvent: weight_logged")
	monitor.LogWeight("42.5")
	fmt.Println("\nEvent: health_event (severity: warning)")
	monitor.LogHealthEvent("limping", "warning")
	fmt.Println("\nEvent: health_event (severity: critical)")
	monitor.LogHealthEvent("PPR", "critical")
}
