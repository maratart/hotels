package business

import (
	"context"
	"errors"
	"fmt"
	"hotels/app/days"
	"hotels/app/logger"
	"hotels/app/model"
	"hotels/app/storage"
	"time"
)

var ErrHotelRoomNotAvailable = errors.New("hotel room is not available for selected dates")
var ErrNotCorrectOrder = errors.New("no correct order")

type OrderService struct {
	log     *logger.Logger
	storage storage.Storage
}

func NewOrderService(log *logger.Logger) OrderService {
	return OrderService{
		log,
		storage.NewStorage(log),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	createdOrder := model.Order{}
	if !s.isValidOrder(order) {
		s.log.Info("invalid order data: %v", order)
		return createdOrder, ErrNotCorrectOrder
	}

	if !s.checkAvailability(order.HotelID, order.RoomID, order.From, order.To) {
		s.log.Info("hotel room is not available for selected dates. %v", order)
		return createdOrder, ErrHotelRoomNotAvailable
	}

	err := s.bookRoom(order)
	if err != nil {
		return createdOrder, fmt.Errorf("failed to book room: %s", err)
	}

	err = s.storage.CreateOrder(ctx, order)
	if err != nil {
		return createdOrder, fmt.Errorf("failed to create order: %s", err)
	}

	createdOrder = order
	return createdOrder, nil
}

func (s *OrderService) isValidOrder(order model.Order) bool {
	return order.HotelID != "" && order.RoomID != "" && order.UserEmail != "" &&
		!order.From.After(order.To) && !order.From.Equal(order.To)
}

func (s *OrderService) checkAvailability(hotelID, roomID string, from, to time.Time) bool {
	daysToBook := days.Between(from, to)
	if daysToBook == nil {
		return false
	}
	availability := s.storage.GetAvailability()

	for _, day := range daysToBook {
		available := false
		for _, availableDay := range availability {
			if availableDay.HotelID == hotelID && availableDay.RoomID == roomID && availableDay.Date.Equal(day) {
				if availableDay.Quota > 0 {
					available = true
					break
				}
			}
		}
		if !available {
			return false
		}
	}
	return true
}

func (s *OrderService) bookRoom(order model.Order) error {
	daysToBook := days.Between(order.From, order.To)
	if daysToBook == nil {
		return ErrHotelRoomNotAvailable
	}
	availability := s.storage.GetAvailability()

	for _, dayToBook := range daysToBook {
		for i, availableDay := range availability {
			if availableDay.HotelID == order.HotelID && availableDay.RoomID == order.RoomID && availableDay.Date.Equal(dayToBook) {
				if availableDay.Quota > 0 {
					availability[i].Quota -= 1
				} else {
					return fmt.Errorf("room not available on %v", dayToBook)
				}
			}
		}
	}

	return s.storage.UpdateAvailability(availability)
}
