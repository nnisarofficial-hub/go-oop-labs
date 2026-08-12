package main

// AnimalMonitor emits events when things happen
type AnimalMonitor struct {
	EventEmitter // embed to get Subscribe/Unsubscribe/Emit for free
	AnimalID     string
}

func (m *AnimalMonitor) LogHealthEvent(condition, severity string) {
	m.Emit(Event{
		Type: "health_event",
		Payload: map[string]string{
			"animal_id": m.AnimalID,
			"condition": condition,
			"severity":  severity,
		},
	})
}

func (m *AnimalMonitor) LogWeight(weightKg string) {
	m.Emit(Event{
		Type:    "weight_logged",
		Payload: map[string]string{"animal_id": m.AnimalID, "weight_kg": weightKg},
	})
}
