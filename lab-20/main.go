package main

import (
	"errors"
	"fmt"
	"time"
)

type AnimalType string

const (
	Doe  AnimalType = "doe"
	Buck AnimalType = "buck"
	Kid  AnimalType = "kid"
)

type Animal struct {
	TagID    string
	Name     string
	Type     AnimalType
	WeightKg float64
}

type HealthRecord struct {
	ID        string
	AnimalTag string
	EventType string // "vaccination", "illness", "checkup"
	Notes     string
	Date      time.Time
	Cost      float64
}

type HealthRepository interface {
	SaveRecord(record HealthRecord) error
	FindByAnimal(tagID string) ([]HealthRecord, error)
	FindAll() ([]HealthRecord, error)
}

type InMemoryHealthRepo struct {
	records []HealthRecord
}

func (i *InMemoryHealthRepo) SaveRecord(record HealthRecord) error {
	i.records = append(i.records, record)
	return nil
}
func (i *InMemoryHealthRepo) FindByAnimal(tagID string) ([]HealthRecord, error) {
	findingAnimal := []HealthRecord{}
	for _, record := range i.records {
		if record.AnimalTag == tagID {
			findingAnimal = append(findingAnimal, record)
		}
	}
	return findingAnimal, nil
}
func (i *InMemoryHealthRepo) FindAll() ([]HealthRecord, error) {
	if len(i.records) == 0 {
		return nil, errors.New("no record found")
	}
	return i.records, nil
}

type CostStrategy interface {
	CalculateCost(event string, baseAmount float64) float64
}

type StandardCost struct{}

func (s StandardCost) CalculateCost(event string, baseAmount float64) float64 {
	if event == "vaccination" {
		return baseAmount * 1.0
	}
	if event == "illness" {
		return baseAmount * 1.5
	}
	if event == "checkup" {
		return 500.0
	}
	return 0.0
}

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

func (c ConsoleLogger) Info(msg string) {
	fmt.Printf("[INFO] [%s] %s\n", c.prefix, msg)
}

func (c ConsoleLogger) Warn(msg string) {
	fmt.Printf("[WARN] [%s] %s\n", c.prefix, msg)
}

func (c ConsoleLogger) Error(msg string) {
	fmt.Printf("[ERROR] [%s] %s\n", c.prefix, msg)
}

type HealthAlertObserver struct {
	logger Logger
}

func (h *HealthAlertObserver) Name() string { return "HealthAlert" }
func (h *HealthAlertObserver) OnEvent(e Event) {
	if e.Type != "illness" {
		return
	}
	fmt.Printf("[%s] 🚨 %s event for %s: %s\n", h.Name(), e.Type, e.Payload["animal_id"], e.Payload["notes"])
}

type Event struct {
	Type    string
	Payload map[string]string
}

type Observer interface {
	OnEvent(event Event)
	Name() string
}

type EventEmitter struct {
	observers []Observer
}

func (e *EventEmitter) Subscribe(o Observer) {
	e.observers = append(e.observers, o)
}

func (e *EventEmitter) Emit(event Event) {
	for _, o := range e.observers {
		o.OnEvent(event)
	}
}

type HealthService struct {
	EventEmitter
	repo     HealthRepository
	logger   Logger
	strategy CostStrategy
}

func NewHealthService(repo HealthRepository, logger Logger, strategy CostStrategy) *HealthService {
	return &HealthService{
		EventEmitter: EventEmitter{},
		repo:         repo,
		logger:       logger,
		strategy:     strategy,
	}
}

func (s *HealthService) RecordEvent(animal Animal, eventType, notes string, baseCost float64) error {
	totalCost := s.strategy.CalculateCost(eventType, baseCost)
	record := HealthRecord{
		ID:        fmt.Sprintf("%s-%d", animal.TagID, time.Now().UnixNano()),
		AnimalTag: animal.TagID,
		EventType: eventType,
		Notes:     notes,
		Date:      time.Now(),
		Cost:      totalCost,
	}
	if err := s.repo.SaveRecord(record); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("recording %s for %s (%s) — cost: PKR %.0f", eventType, animal.TagID, animal.Name, totalCost))
	s.Emit(Event{
		Type: eventType,
		Payload: map[string]string{
			"animal_id": animal.TagID,
			"name":      animal.Name,
			"notes":     notes,
			"cost":      fmt.Sprintf("%.0f", totalCost),
		},
	})
	return nil
}

func (s *HealthService) GetAnimalHistory(tagID string) ([]HealthRecord, error) {
	animalHistory, err := s.repo.FindByAnimal(tagID)
	if err != nil {
		return nil, err
	}
	return animalHistory, nil
}

func (s *HealthService) GetTotalCost(from, to time.Time) (float64, error) {
	records, err := s.repo.FindAll()
	if err != nil {
		return 0.0, err
	}
	var total float64
	for _, r := range records {
		if !r.Date.Before(from) && !r.Date.After(to) {
			total += r.Cost
		}
	}
	return total, nil
}
func main() {
	repo := &InMemoryHealthRepo{}
	logger := NewConsoleLogger("health")
	strategy := StandardCost{}
	svc := NewHealthService(repo, logger, strategy)
	svc.Subscribe(&HealthAlertObserver{logger: logger})
	doe := Animal{TagID: "D-001", Name: "Rani", Type: Doe, WeightKg: 38}
	svc.RecordEvent(doe, "vaccination", "PPR vaccine — annual", 1500)
	svc.RecordEvent(doe, "checkup", "routine monthly checkup", 0)
	svc.RecordEvent(doe, "illness", "limping — possible injury", 2000)
	history, _ := svc.GetAnimalHistory("D-001")
	fmt.Printf("\nHealth history for %s (%s):\n", doe.Name, doe.TagID)
	for _, r := range history {
		fmt.Printf("  [%s] %s — PKR %.0f\n", r.EventType, r.Notes, r.Cost)
	}
	total, _ := svc.GetTotalCost(time.Now().Add(-7*24*time.Hour), time.Now())
	fmt.Printf("\nTotal health cost this week: PKR %.0f\n", total)
}
