package main

import (
	"errors"
	"fmt"
	"sort"
)

type ProductService struct {
	repo   ProductRepository
	logger Logger
}

func NewProductService(repo ProductRepository, logger Logger) *ProductService {
	return &ProductService{repo: repo, logger: logger}
}

func (s *ProductService) AddProduct(id, name, category string, price float64, stock int) error {
	if name == "" {
		return errors.New("no name")
	}
	if price < 0 {
		return errors.New("price cannot be less than zero")
	}
	if stock < 0 {
		return errors.New("stock cannot be less than zero")
	}
	product := &Product{
		ID:       id,
		Name:     name,
		Price:    price,
		Stock:    stock,
		Category: category,
	}
	s.repo.Save(product)
	return nil
}

func (s *ProductService) GetCatalog() ([]*Product, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	sort.Slice(products, func(i, j int) bool {
		if products[i].Category != products[j].Category {
			return products[i].Category < products[j].Category
		}
		return products[i].Name < products[j].Name
	})
	return products, nil
}

func (s *ProductService) UpdateStock(id string, delta int) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}
	newStock := product.Stock + delta
	if newStock < 0 {
		return fmt.Errorf("insufficient stock (have %d, requested %d)", product.Stock, -delta)
	}
	product.Stock = newStock
	s.logger.Info(fmt.Sprintf("Updated stock for %s: %d -> %d", id, product.Stock-delta, newStock))
	return s.repo.Save(product)
}

func (s *ProductService) GetLowStockAlerts(threshold int) ([]*Product, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var lowStockProduct []*Product
	for _, product := range products {
		if product.Stock <= threshold {
			lowStockProduct = append(lowStockProduct, product)
		}
	}
	return lowStockProduct, nil
}
