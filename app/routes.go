package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func newMux(h *HandlerService) http.Handler {
	r := chi.NewRouter()
	r.Post("/orders", h.createOrderHandler)
	// other routes...

	return r
}
