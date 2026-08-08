package main

import "fmt"

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type OrderService struct {
	logger Logger
	orders map[int]string
}

func NewOrderService(logger Logger) *OrderService {
	return &OrderService{logger: logger, orders: map[int]string{}}
}

func (s *OrderService) PlaceOrder(id int, item string) {
	s.logger.Info(fmt.Sprintf("placing order %d for item: %s", id, item))
	s.orders[id] = item
	s.logger.Info(fmt.Sprintf("order %d placed successfully", id))
}

func (s *OrderService) CancelOrder(id int) error {
	if _, exists := s.orders[id]; !exists {
		s.logger.Warn(fmt.Sprintf("attempted to cancel non-existent order %d", id))
		return fmt.Errorf("order %d not found", id)
	}
	delete(s.orders, id)
	s.logger.Info(fmt.Sprintf("order %d cancelled", id))
	return nil
}
