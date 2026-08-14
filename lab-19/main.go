package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type ReportData struct {
	Title   string
	From    time.Time
	To      time.Time
	Records []map[string]string
}

type ReportGenerator interface {
	Name() string
	Description() string
	Generate(data ReportData) (string, error)
}

type ReportRegistry struct {
	generators map[string]ReportGenerator
}

func NewReportRegistry() *ReportRegistry {
	return &ReportRegistry{
		generators: make(map[string]ReportGenerator),
	}
}

func (r *ReportRegistry) Register(g ReportGenerator) error {
	_, exists := r.generators[g.Name()]
	if exists {
		return errors.New("error: key already exists in the map")
	}
	r.generators[g.Name()] = g
	return nil
}

func (r *ReportRegistry) Get(name string) (ReportGenerator, error) {
	value, exists := r.generators[name]
	if !exists {
		return nil, errors.New("error: report generator 'nonexistent' not found")
	}
	return value, nil
}

func (r *ReportRegistry) List() []string {
	names := []string{}
	for name := range r.generators {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

type AnimalCountReport struct{}

func (a *AnimalCountReport) Name() string        { return "animal-count" }
func (a *AnimalCountReport) Description() string { return "Count of animals by type" }
func (a *AnimalCountReport) Generate(data ReportData) (string, error) {
	if len(data.Records) == 0 {
		return "", errors.New("no records to generate report from")
	}
	counts := make(map[string]int)
	for _, record := range data.Records {
		animalType := record["type"]
		counts[animalType]++
	}
	result := "Animal Count Report\n"
	for animalType, count := range counts {
		result += fmt.Sprintf("%s: %d\n", animalType, count)
	}
	return result, nil
}

type HealthSummaryReport struct{}

func (h *HealthSummaryReport) Name() string        { return "health-summary" }
func (h *HealthSummaryReport) Description() string { return "Health events summary" }
func (h *HealthSummaryReport) Generate(data ReportData) (string, error) {
	if len(data.Records) == 0 {
		return "", errors.New("no records to generate report from")
	}
	counts := make(map[string]int)
	for _, record := range data.Records {
		eventType := record["event_type"]
		counts[eventType]++
	}
	result := "Health Summary Report\n"
	for animalType, count := range counts {
		result += fmt.Sprintf("%s: %d\n", animalType, count)
	}
	return result, nil
}

type FinancialReport struct{}

func (f *FinancialReport) Name() string        { return "financial" }
func (f *FinancialReport) Description() string { return "Income and expense summary" }
func (f *FinancialReport) Generate(data ReportData) (string, error) {
	if len(data.Records) == 0 {
		return "", errors.New("no records to generate report from")
	}
	var totalIncome, totalExpense float64
	for _, record := range data.Records {
		amount, _ := strconv.ParseFloat(record["amount"], 64)
		if record["type"] == "income" {
			totalIncome += amount
		} else if record["type"] == "expense" {
			totalExpense += amount
		}
	}
	net := totalIncome - totalExpense
	dateRange := data.From.Format("Jan 02") + " – " + data.To.Format("Jan 02")

	result := fmt.Sprintf(
		"Financial Report: %s\n  Total income:   PKR %s\n  Total expenses: PKR %s\n  Net profit:     PKR %s\n",
		dateRange,
		formatWithCommas(totalIncome),
		formatWithCommas(totalExpense),
		formatWithCommas(net),
	)
	return result, nil
}

func formatWithCommas(n float64) string {
	str := fmt.Sprintf("%.0f", n)
	var result []byte
	for i, digit := range []byte(str) {
		if i != 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, digit)
	}
	return string(result)
}

func main() {
	registry := NewReportRegistry()
	registry.Register(&AnimalCountReport{})
	registry.Register(&HealthSummaryReport{})
	registry.Register(&FinancialReport{})

	fmt.Println("Available reports:", registry.List())
	fmt.Println()
	data := ReportData{
		Title: "Farm Financial Report",
		From:  time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
		Records: []map[string]string{
			{"type": "income", "amount": "300000"},
			{"type": "income", "amount": "150000"},
			{"type": "expense", "amount": "120000"},
			{"type": "expense", "amount": "65000"},
		},
	}
	gen, _ := registry.Get("financial")
	output, _ := gen.Generate(data)
	fmt.Println(output)
	_, err := registry.Get("nonexistent")
	fmt.Println(err)
}
