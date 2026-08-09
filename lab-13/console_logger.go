package main

import "fmt"

type ConsoleLogger struct{}

func (c *ConsoleLogger) Info(msg string)  { fmt.Println("[INFO]", msg) }
func (c *ConsoleLogger) Warn(msg string)  { fmt.Println("[WARN]", msg) }
func (c *ConsoleLogger) Error(msg string) { fmt.Println("[ERROR]", msg) }
func printCatalog(products []*Product) {
	fmt.Printf("Catalog (%d products):\n", len(products))

	for _, p := range products {
		fmt.Printf("  [%-11s] %-15s — PKR %8.0f | Stock: %d\n",
			p.Category, p.Name, p.Price, p.Stock)
	}
}
