package main

import "fmt"

type EventEmitter struct {
	observers []Observer
}

func (e *EventEmitter) Subscribe(o Observer) {
	e.observers = append(e.observers, o)
}

func (e *EventEmitter) Unsubscribe(name string) {
	if len(e.observers) == 0 {
		fmt.Println("there is no subscribers")
		return
	}
	var result []Observer
	for _, o := range e.observers {
		if o.Name() != name {
			result = append(result, o)
		}
	}
	e.observers = result
}

func (e *EventEmitter) Emit(event Event) {
	for _, o := range e.observers {
		o.OnEvent(event)
	}
}
