package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
// @Summary      Send single notification
// @Description  Creates a new notification and sends it.
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        notification  body      domain.Notification  true  "Notification details"
// @Success      202  {object}  map[string]string "Returns notification ID and status"
// @Router       /notifications [post]
func (h *NotificationHandler) Send(w http.ResponseWriter, r *http.Request) {
	var n domain.Notification

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Trigger usecase
	if err := h.usecase.Execute(r.Context(), &n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"id": n.ID.String(), "status": string(n.Status)})
}

// SendBatch handles multiple notification requests in a single call
// SendBatch handles multiple notification requests
// @Summary      Send batch notifications
// @Description  Creates multiple notifications in a single request.
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        notifications body []domain.Notification true "List of notifications"
// @Success      202 {object} map[string]interface{}
// @Router       /notifications/batch [post]
func (h *NotificationHandler) SendBatch(w http.ResponseWriter, r *http.Request) {
	var ns []domain.Notification

	if err := json.NewDecoder(r.Body).Decode(&ns); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	batchID, err := h.usecase.ExecuteBatch(r.Context(), ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"batch_id": batchID,
		"count":    len(ns),
		"status":   "accepted",
	})
}

// List handles notification filtering and pagination
// List handles notification filtering and pagination
// @Summary      List notifications
// @Description  Retrieves notifications with filtering and pagination.
// @Tags         Notifications
// @Produce      json
// @Param        status    query     string  false  "Filter by status"
// @Param        channel   query     string  false  "Filter by channel"
// @Param        limit     query     int     false  "Limit for pagination" default(20)
// @Param        offset    query     int     false  "Offset for pagination"
// @Success      200 {array} domain.Notification
// @Router       /notifications [get]
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(query.Get("offset"))

	filter := domain.NotificationFilter{
		Status:   query.Get("status"),
		Channel:  query.Get("channel"),
		Priority: query.Get("priority"),
		Limit:    limit,
		Offset:   offset,
	}

	if bID := query.Get("batch_id"); bID != "" {
		filter.BatchID = &bID
	}

	if start := query.Get("start_date"); start != "" {
		t, _ := time.Parse(time.RFC3339, start)
		filter.StartDate = &t
	}
	if end := query.Get("end_date"); end != "" {
		t, _ := time.Parse(time.RFC3339, end)
		filter.EndDate = &t
	}

	notifications, err := h.usecase.List(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// Cancel handles bulk cancellation using path parameters
// Cancel handles bulk cancellation
// @Summary      Cancel batch
// @Description  Cancels all pending notifications in a batch.
// @Tags         Notifications
// @Param        batchId   path      string  true  "Batch ID"
// @Success      204 "No Content"
// @Router       /notifications/batch/{batchId}/cancel [patch]
func (h *NotificationHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	batchID := r.PathValue("batchId")
	if batchID == "" {
		http.Error(w, "batchId is required", http.StatusBadRequest)
		return
	}

	if err := h.usecase.CancelBatch(r.Context(), batchID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Metrics handles system observability
// @Summary      System Metrics
// @Description  Get current status of queues and system health.
// @Tags         Observability
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /metrics [get]
func (h *NotificationHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.usecase.GetSystemMetrics(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch metrics", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"queues":    stats,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
