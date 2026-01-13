package http

import (
	"net/http"

	"karhub-beer-machine/internal/interfaces/http/handlers"
)

// RegisterRoutes sets up the HTTP routes for the beer machine application.
func RegisterRoutes(mux *http.ServeMux, h *handlers.BeerHandler) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// CRUD
	mux.HandleFunc("/beer-styles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/beer-styles/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Core business
	mux.HandleFunc("/beer-styles/best", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.FindBest(w, r)
	})
}
