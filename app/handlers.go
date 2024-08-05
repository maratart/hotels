package app

import (
	"encoding/json"
	"errors"
	"hotels/app/business"
	"hotels/app/logger"
	"hotels/app/model"
	"net/http"
)

type HandlerService struct {
	log *logger.Logger
	srv business.OrderService
}

func NewHandlers(log *logger.Logger) *HandlerService {
	return &HandlerService{
		log,
		business.NewOrderService(log),
	}
}

func (h *HandlerService) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		return
	}

	ctx := r.Context()
	createdOrder, err := h.srv.CreateOrder(ctx, newOrder)
	if err != nil {
		if errors.Is(err, business.ErrHotelRoomNotAvailable) {
			http.Error(w, "hotel room is not available for selected dates", http.StatusBadRequest)
			return
		}
		if errors.Is(err, business.ErrNotCorrectOrder) {
			http.Error(w, "not correct order", http.StatusBadRequest)
			return
		}
		h.log.Errorf("error on order creation: %s", err)
		http.Error(w, "error on order creation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(createdOrder)
	if err != nil {
		h.log.Errorf("error on order encoding: %s", err)
		http.Error(w, "error on order encoding", http.StatusInternalServerError)
		return
	}

	h.log.Info("order successfully created: %v", newOrder)
}
