package main

import (
	"errors"
	"fmt"
)

type Color struct {
	R, G, B uint8 // 0–255
}

// String() → "rgb(255, 128, 0)"
func (c Color) String() string {
	return fmt.Sprintf("Color: rgb(%d,%d,%d)", c.R, c.G, c.B)
}

// Also implement: Hex() string → "#FF8000"
func (c Color) Hex() string {
	return fmt.Sprintf("HEX:   #%02X%02X%02X", c.R, c.G, c.B)
}

type Temperature struct {
	Celsius float64
}

// String() → "37.0°C (98.6°F)"
func (t Temperature) String() string {
	fahrenheit := (t.Celsius * 9 / 5) + 32
	return fmt.Sprintf("Temperature: %.1f°C (%.1f°F)", t.Celsius, fahrenheit)
}

// Also implement: IsFever() bool — true if > 37.5°C
func (t Temperature) IsFever() bool {
	return t.Celsius > 37.5
}

type Money struct {
	Amount   float64
	Currency string
}

// String() → "PKR 1,250.00" or "USD 9.99"
func (m Money) String() string {
	return fmt.Sprintf("Price: %s %.2f", m.Currency, m.Amount)
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("Error: cannot add PKR and USD")
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

func (m Money) IsZero() bool {
	return m.Amount == 0.0
}

// Also implement: Add(other Money) (Money, error) — error if currencies differ
//                 IsZero() bool

func main() {
	hex := Color{
		R: 255,
		G: 128,
		B: 0,
	}
	fmt.Println(hex)
	fmt.Println(hex.Hex())

	temp := Temperature{
		Celsius: 38.5,
	}
	fmt.Println()
	fmt.Println(temp)
	fmt.Printf("Is fever: %t\n", temp.IsFever())

	m1 := Money{Amount: 1250.00, Currency: "PKR"}
	m2 := Money{Amount: 1250.00, Currency: "PKR"}
	m3 := Money{Amount: 10.00, Currency: "USD"}
	fmt.Printf("\nPrice: %s\n", m1)
	total, _ := m1.Add(m2)
	fmt.Printf("Total: %s\n\n", total)
	_, err := m1.Add(m3)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
