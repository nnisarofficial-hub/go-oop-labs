package main

import (
	"fmt"
	"strconv"
)

type PricingStrategy interface {
	Apply(basePrice float64) float64
	Description() string
}

type NoDiscount struct{}

func (n NoDiscount) Apply(price float64) float64 { return price }
func (n NoDiscount) Description() string         { return "No discount" }

type PercentageDiscount struct {
	Percent float64
}

func (p PercentageDiscount) Apply(price float64) float64 {
	discount := price * (p.Percent / 100)
	finalPrice := price - discount
	return finalPrice
}
func (p PercentageDiscount) Description() string {
	return fmt.Sprintf("%f discount", p.Percent)
}

type QurbaniSeasonDiscount struct{}

func (q QurbaniSeasonDiscount) Apply(price float64) float64 {
	discount := price * (20.0 / 100.0)
	finalPrice := price - discount - 500
	return finalPrice
}
func (q QurbaniSeasonDiscount) Description() string {
	return "Qurbani Season (20% off + free delivery)"
}

type BulkDiscount struct {
	Quantity  int
	Threshold int
	Percent   float64
}

func (b BulkDiscount) Apply(price float64) float64 {
	if b.Quantity > b.Threshold {
		discount := price * (b.Percent / 100.0)
		finalPrice := price - discount
		return finalPrice
	} else {
		return price
	}
}
func (b BulkDiscount) Description() string {
	return fmt.Sprintf("%.0f%% bulk discount (10+ units)", b.Percent)
}

type Order struct {
	ID        string
	Item      string
	Quantity  int
	UnitPrice float64
	strategy  PricingStrategy
}

func NewOrder(id, item string, qty int, unitPrice float64, strategy PricingStrategy) *Order {
	return &Order{
		ID:        id,
		Item:      item,
		Quantity:  qty,
		UnitPrice: unitPrice,
		strategy:  strategy,
	}
}

func (o *Order) Total() float64 {
	return o.strategy.Apply(o.UnitPrice * float64(o.Quantity))
}

func (o *Order) PrintReceipt() {
	fmt.Printf("Order: %s\n", o.ID)
	fmt.Printf("  %-10s%s\n", "Item:", o.Item)
	fmt.Printf("  %-10s%d x PKR %s\n", "Qty:", o.Quantity, formatWithCommas(o.UnitPrice))
	fmt.Printf("  %-10s%s\n", "Pricing:", o.strategy.Description())
	fmt.Printf("  %-10sPKR %s", "Total:", formatWithCommas(o.Total()))
}

func formatWithCommas(n float64) string {
	intPart := int64(n)
	dec := int64((n - float64(intPart)) * 100)

	s := strconv.FormatInt(intPart, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return fmt.Sprintf("%s.%02d", s, dec)
}

func main() {
	strategy1 := NoDiscount{}
	strategy2 := QurbaniSeasonDiscount{}
	strategy3 := BulkDiscount{Quantity: 15, Threshold: 10, Percent: 15}
	order1 := NewOrder("ORD-001", "Barbari Buck", 1, 50000, strategy1)
	order1.PrintReceipt()
	fmt.Println()
	order2 := NewOrder("ORD-002", "Barbari Buck", 1, 50000, strategy2)
	order2.PrintReceipt()
	fmt.Println()
	order3 := NewOrder("ORD-003", "Barbari Buck", 15, 800, strategy3)
	order3.PrintReceipt()
	fmt.Println()
}
