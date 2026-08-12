package main

type Event struct {
	Type    string // "status_changed", "weight_logged", "health_event"
	Payload map[string]string
}
