package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/ahmetcanc/notify-one/internal/usecase"
)

type NotificationHandler struct {
	usecase *usecase.NotificationUsecase
}

func NewNotificationHandler(u *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{usecase: u}
}

// Send handles single notification requests
func (h *NotificationHandler) Send(w http.ResponseWriter, r *http.Request) {
	var n domain.Notification

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Trigger usecase
	if err := h.usecase.Execute(r.Context(), &n); err != nil {
		http.Error(w, "failed to process notification", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"id": n.ID.String(), "status": string(n.Status)})
}

// SendBatch handles multiple notification requests in a single call
func (h *NotificationHandler) SendBatch(w http.ResponseWriter, r *http.Request) {
	var ns []domain.Notification

	if err := json.NewDecoder(r.Body).Decode(&ns); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.ExecuteBatch(r.Context(), ns); err != nil {
		http.Error(w, "failed to process batch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "batch processed", "count": fmt.Sprintf("%d", len(ns))})
}
