package main

import "fmt"

type HealthTrackerObserver struct{}

func (h *HealthTrackerObserver) Name() string { return "HealthTracker" }
func (h *HealthTrackerObserver) OnEvent(event Event) {
	if event.Type == "health_event" {
		fmt.Printf("  [%s] recording health event for %s: %s (%s)\n", h.Name(), event.Payload["animal_id"], event.Payload["condition"], event.Payload["severity"])
	} else {
		fmt.Printf("  [%s] ignoring non-health event\n", h.Name())
	}
}
