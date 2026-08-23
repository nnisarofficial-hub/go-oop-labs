package main

import "errors"

type InMemoryProductRepo struct {
	products map[string]*Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{products: make(map[string]*Product)}
}
func (r *InMemoryProductRepo) FindByID(id string) (*Product, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	product := r.products[id]
	return product, nil
}
func (r *InMemoryProductRepo) FindAll() ([]*Product, error) {
	var allProducts []*Product
	for _, product := range r.products {
		allProducts = append(allProducts, product)
	}
	return allProducts, nil
}
func (r *InMemoryProductRepo) FindByCategory(category string) ([]*Product, error) {
	if category == "" {
		return nil, errors.New("category cannot be empty")
	}
	var filteredProducts []*Product
	for _, product := range r.products {
		if product.Category == category {
			filteredProducts = append(filteredProducts, product)
		}
	}
	return filteredProducts, nil
}
func (r *InMemoryProductRepo) Save(p *Product) error {
	if p == nil {
		return errors.New("cannot add empty product")
	}
	if p.ID == "" {
		return errors.New("product id cannot be empty")
	}
	r.products[p.ID] = p
	return nil
}
func (r *InMemoryProductRepo) Delete(id string) error {
	return nil
}
