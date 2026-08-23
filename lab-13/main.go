package main

import "fmt"

type Product struct {
	ID       string
	Name     string
	Price    float64
	Stock    int
	Category string
}

type ProductRepository interface {
	FindByID(id string) (*Product, error)
	FindAll() ([]*Product, error)
	FindByCategory(category string) ([]*Product, error)
	Save(p *Product) error
	Delete(id string) error
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func printCatalog(products []*Product) {
	fmt.Printf("Catalog (%d products):\n", len(products))

	for _, p := range products {
		fmt.Printf("  [%-11s] %-15s — PKR %8.0f | Stock: %d\n",
			p.Category, p.Name, p.Price, p.Stock)
	}
}

func main() {
	repo := NewInMemoryProductRepo()
	logger := &ConsoleLogger{}
	service := NewProductService(repo, logger)

	service.AddProduct("p1", "Laptop Pro 15", "electronics", 250000, 5)
	service.AddProduct("p2", "Wireless Mouse", "electronics", 3500, 50)
	service.AddProduct("p3", "Barseem Hay", "feed", 800, 3)
	service.AddProduct("p4", "Ivermectin", "medicine", 1200, 2)

	catalog, err := service.GetCatalog()
	if err != nil {
		fmt.Println("Error fetching catalog:", err)
		return
	}

	printCatalog(catalog)
	fmt.Println()
	threshold := 5
	lowStock, err := service.GetLowStockAlerts(threshold)
	if err != nil {
		fmt.Println("Error fetching low stock alerts:", err)
		return
	}

	fmt.Printf("\nLow stock (threshold: %d):\n", threshold)
	for _, p := range lowStock {
		fmt.Printf("  %s (%d remaining)\n", p.Name, p.Stock)
	}

	fmt.Println("\nAfter selling 4 units of Barseem Hay:")
	err = service.UpdateStock("p3", -4)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
