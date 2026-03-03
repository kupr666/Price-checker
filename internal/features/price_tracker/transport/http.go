package transport

import (
	"context"
	"errors"
	"encoding/json"
	"net/http"

	"price_checker/internal/core/domains"
	
	"go.uber.org/zap"
)

type Service interface {
	 CreateItem(ctx context.Context, item domains.Item) (domains.Item, error)
	 ListItems(ctx context.Context) ([]domains.Item, error)
}

type Handler struct {
	svc Service
	logger *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		svc: service,
		logger: logger,
	}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var item domains.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		h.logger.Warn("failed to decode requestbody", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	created, err := h.svc.CreateItem(r.Context(), item)
	if err != nil {
		if errors.Is(err, domains.ErrUrlExists) {
			h.logger.Info("Attempting to add existing url", zap.String("url", item.URL))
			http.Error(w, err.Error(), http.StatusConflict)
			return 
		}
		
		h.logger.Error("failed to create item", zap.Error(err), zap.String("url", item.URL))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("New item was created", zap.Int64("id", created.ID), zap.String("url", created.URL))
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListItems(r.Context())
	if err != nil {
		h.logger.Error("failed to get all items from list", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(items)
}