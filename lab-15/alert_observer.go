package main

import "fmt"

type AlertObserver struct {
	threshold string // only alert on "critical" events
}

func (a *AlertObserver) Name() string { return "AlertSystem" }
func (a *AlertObserver) OnEvent(event Event) {
	// only alert if payload["severity"] == "critical"
	if event.Type != "health_event" {
		return
	}
	if event.Payload["severity"] == "critical" {
		fmt.Printf("  [%s] 🚨 CRITICAL ALERT: animal %s has %s\n", a.Name(), event.Payload["animal_id"], event.Payload["condition"])
	} else {
		fmt.Printf("  [%s] no alert — severity is warning, not critical\n", a.Name())
	}
}
