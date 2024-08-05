package storage

import (
	"context"
	"hotels/app/days"
	"hotels/app/logger"
	"hotels/app/model"
	"sync"
)

type Storage interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrders() []model.Order
	GetAvailability() []model.RoomAvailability
	UpdateAvailability(availability []model.RoomAvailability) error
}

type storageImpl struct {
	log          *logger.Logger
	orders       []model.Order
	availability []model.RoomAvailability
	mu           sync.RWMutex
}

func NewStorage(log *logger.Logger) Storage {
	var orders []model.Order
	availability := []model.RoomAvailability{
		{"reddison", "lux", days.Date(2024, 1, 1), 1},
		{"reddison", "lux", days.Date(2024, 1, 2), 1},
		{"reddison", "lux", days.Date(2024, 1, 3), 1},
		{"reddison", "lux", days.Date(2024, 1, 4), 1},
		{"reddison", "lux", days.Date(2024, 1, 5), 0},
	}

	return &storageImpl{
		log,
		orders,
		availability,
		sync.RWMutex{},
	}
}

func (s *storageImpl) CreateOrder(_ context.Context, order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders = append(s.orders, order)
	s.log.Info("order successfully created: %v", order)
	return nil
}

func (s *storageImpl) GetOrders() []model.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.orders
}

func (s *storageImpl) GetAvailability() []model.RoomAvailability {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.availability
}

func (s *storageImpl) UpdateAvailability(availability []model.RoomAvailability) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.availability = availability
	return nil
}
