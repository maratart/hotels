package business

import (
	"context"
	"errors"
	"hotels/app/days"
	"hotels/app/logger"
	"hotels/app/model"
	"testing"
	"time"
)

// MockStorage is a mock implementation of the storage.Storage interface.
type MockStorage struct {
	createOrderFunc        func(ctx context.Context, order model.Order) error
	getAvailabilityFunc    func() []model.RoomAvailability
	updateAvailabilityFunc func([]model.RoomAvailability) error
}

func (m *MockStorage) CreateOrder(ctx context.Context, order model.Order) error {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, order)
	}
	return nil
}

func (m *MockStorage) GetOrders() []model.Order {
	return nil // Not used in tests
}

func (m *MockStorage) GetAvailability() []model.RoomAvailability {
	if m.getAvailabilityFunc != nil {
		return m.getAvailabilityFunc()
	}
	return nil
}

func (m *MockStorage) UpdateAvailability(availability []model.RoomAvailability) error {
	if m.updateAvailabilityFunc != nil {
		return m.updateAvailabilityFunc(availability)
	}
	return nil
}

func TestCreateOrder(t *testing.T) {
	mockStorage := &MockStorage{}
	orderService := OrderService{log: logger.NewLogger(), storage: mockStorage}

	tests := []struct {
		name          string
		order         model.Order
		setupMock     func()
		expectedOrder model.Order
		expectedErr   error
	}{
		{
			name: "invalid order",
			order: model.Order{
				HotelID:   "",
				RoomID:    "101",
				UserEmail: "test@example.com",
				From:      time.Now(),
				To:        time.Now().Add(2 * time.Hour),
			},
			setupMock:   func() {},
			expectedErr: ErrNotCorrectOrder,
		},
		{
			name: "room not available",
			order: model.Order{
				HotelID:   "reddison",
				RoomID:    "lux",
				UserEmail: "test@example.com",
				From:      days.Date(2024, 1, 5),
				To:        days.Date(2024, 1, 6),
			},
			setupMock: func() {
				mockStorage.getAvailabilityFunc = func() []model.RoomAvailability {
					return []model.RoomAvailability{
						{"reddison", "lux", days.Date(2024, 1, 5), 0},
					}
				}
			},
			expectedErr: ErrHotelRoomNotAvailable,
		},
		{
			name: "successful order creation",
			order: model.Order{
				HotelID:   "reddison",
				RoomID:    "lux",
				UserEmail: "test@example.com",
				From:      days.Date(2024, 1, 1),
				To:        days.Date(2024, 1, 3),
			},
			setupMock: func() {
				mockStorage.getAvailabilityFunc = func() []model.RoomAvailability {
					return []model.RoomAvailability{
						{"reddison", "lux", days.Date(2024, 1, 1), 1},
						{"reddison", "lux", days.Date(2024, 1, 2), 1},
					}
				}
				mockStorage.createOrderFunc = func(ctx context.Context, order model.Order) error {
					return nil
				}
				mockStorage.updateAvailabilityFunc = func([]model.RoomAvailability) error {
					return nil
				}
			},
			expectedOrder: model.Order{
				HotelID:   "reddison",
				RoomID:    "lux",
				UserEmail: "test@example.com",
				From:      days.Date(2024, 1, 1),
				To:        days.Date(2024, 1, 3),
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			createdOrder, err := orderService.CreateOrder(context.Background(), tt.order)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if createdOrder != tt.expectedOrder {
				t.Fatalf("expected order %v, got %v", tt.expectedOrder, createdOrder)
			}
		})
	}
}
