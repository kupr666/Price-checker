package transport

import (
	"errors"
	"encoding/json"
	"net/http"
	"price_checker/internal/core/domains"
)

type Service interface {
	 CreateItem(domains.Item) (domains.Item, error)
	 ListItems() ([]domains.Item, error)
}

type Handler struct {
	svc Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		svc: service,
	}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var item domains.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	created, err := h.svc.CreateItem(item)
	if err != nil {
		if errors.Is(err, domains.ErrUrlExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return 
		}
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "applications/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListItems()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(items)
}