package main

import (
	"errors"
	"fmt"
)

type Vehicle struct {
	Make         string
	Model        string
	Year         int
	Mileage      float64 // km
	FuelCapacity float64 // liters
}

func (v Vehicle) Info() string {
	return fmt.Sprintf("%d %s %s (%.0f km)", v.Year, v.Make, v.Model, v.Mileage)
}

func (v *Vehicle) Drive(km float64) {
	v.Mileage += km
}

type Truck struct {
	Vehicle
	PayloadCapacity float64 // tonnes
	CurrentLoad     float64
}

func NewTruck(make, model string, year int, payload float64) *Truck {
	return &Truck{
		Vehicle: Vehicle{
			Make:  make,
			Model: model,
			Year:  year,
		},
		PayloadCapacity: payload,
		CurrentLoad:     0.0,
	}
}

func (t *Truck) LoadCargo(tonnes float64) error {
	remaining := t.PayloadCapacity - t.CurrentLoad
	if tonnes > remaining {
		return fmt.Errorf("cannot load %.1ft, only %.1ft capacity remaining", tonnes, remaining)
	}
	t.CurrentLoad += tonnes
	return nil
}

func (t *Truck) Info() string {
	return fmt.Sprintf("%s | Load: %.1f/%.1ft", t.Vehicle.Info(), t.CurrentLoad, t.PayloadCapacity)
}

type ElectricCar struct {
	Vehicle
	BatteryCapacity float64 // kwh
	ChargeLevel     float64 // 0.0 to 1.0 (percentage as decimal)
}

func NewElectricCar(make, model string, year int, battery float64) *ElectricCar {
	return &ElectricCar{
		Vehicle: Vehicle{
			Make:  make,
			Model: model,
			Year:  year,
		},
		BatteryCapacity: battery,
		ChargeLevel:     0.0,
	}
}

func (e *ElectricCar) Charge(toLevel float64) error {
	if toLevel > 1.0 {
		return errors.New("requested charge level exceeds maximum capacity")
	}
	if toLevel < e.ChargeLevel {
		return errors.New("requested charge level is below current charge level")
	}
	e.ChargeLevel = toLevel
	return nil
}

func (e *ElectricCar) RangeRemaining() float64 {
	return e.BatteryCapacity * e.ChargeLevel * 6
}

func (e *ElectricCar) Info() string {
	return fmt.Sprintf("%s | Battery: %.0f | Range: ~%.0f km", e.Vehicle.Info(), e.ChargeLevel, e.RangeRemaining())
}

func main() {
	vehicle := Vehicle{Make: "Tesla", Model: "3", Year: 2023, Mileage: 0}
	vehicle3 := Vehicle{Make: "Tesla", Model: "3", Year: 2023, Mileage: 250}
	vehicle1 := Vehicle{Make: "Toyota", Model: "Hilux", Year: 2020, Mileage: 0}
	vehicle2 := Vehicle{Make: "Toyota", Model: "Hilux", Year: 2020, Mileage: 500}

	tTruck := Truck{Vehicle: vehicle1, PayloadCapacity: 5.0, CurrentLoad: 0.0}
	fmt.Println(tTruck.Info())

	cargoLoad := 3.0
	err := tTruck.LoadCargo(cargoLoad)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("After loading %.f tonnes: \n", cargoLoad)
	}

	tTruck.Vehicle = vehicle2
	fmt.Println(tTruck.Info())

	fmt.Println()

	eCar := ElectricCar{Vehicle: vehicle, BatteryCapacity: 60, ChargeLevel: 0.2}
	eCar1 := ElectricCar{Vehicle: vehicle3, BatteryCapacity: 60, ChargeLevel: 0.9}

	fmt.Println(eCar.Info())

	chargeTarget := 0.9
	err = eCar1.Charge(chargeTarget)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("After charging to %.0f%%\n", chargeTarget*100)
	}
	fmt.Println(eCar1.Info())

	fmt.Println()

	err = tTruck.LoadCargo(4.0)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
