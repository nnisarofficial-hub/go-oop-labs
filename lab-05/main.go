package main

import (
	"fmt"
)

type Color struct {
	R, G, B uint8 // 0–255
}

func (c Color) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

func (c Color) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

type Temperature struct {
	Celsius float64
}

func (t Temperature) String() string {
	fahrenheit := (t.Celsius * 9 / 5) + 32
	return fmt.Sprintf("%.1f°C (%.1f°F)", t.Celsius, fahrenheit)
}

func (t Temperature) IsFever() bool {
	return t.Celsius > 37.5
}

type Money struct {
	Amount   float64
	Currency string
}

func (m Money) String() string {
	return fmt.Sprintf("%s %.2f", m.Currency, m.Amount)
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("cannot add %s and %s", m.Currency, other.Currency)
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

func (m Money) IsZero() bool {
	return m.Amount == 0.0
}

func main() {
	color := Color{
		R: 255,
		G: 128,
		B: 0,
	}
	fmt.Printf("Color: %s\n", color.String())
	fmt.Printf("Hex:   %s\n", color.Hex())
	temp := Temperature{
		Celsius: 38.5,
	}
	fmt.Println()
	fmt.Printf("Temperature: %s\n", temp.String())
	fmt.Printf("Is fever: %t\n", temp.IsFever())
	m1 := Money{Amount: 1250.00, Currency: "PKR"}
	m2 := Money{Amount: 1250.00, Currency: "PKR"}
	m3 := Money{Amount: 10.00, Currency: "USD"}
	fmt.Printf("\nPrice: %s\n", m1)
	total, err := m1.Add(m2)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("Total: %s\n\n", total)
	_, err = m1.Add(m3)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
