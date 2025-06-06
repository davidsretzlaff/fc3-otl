package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"

	"payments-subscription/internal/common/logging"

	"github.com/gorilla/mux"
)

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error         string `json:"error"`
	Message       string `json:"message"`
	CorrelationID string `json:"correlation_id"`
	StatusCode    int    `json:"status_code"`
}

// SuccessResponse representa uma resposta de sucesso padronizada
type SuccessResponse struct {
	Data          interface{} `json:"data"`
	Message       string      `json:"message,omitempty"`
	CorrelationID string      `json:"correlation_id"`
}

// handler gerencia as requisições HTTP para Subscription
type handler struct {
	service SubscriptionServiceInterface
}

// NewSubscriptionHandler cria uma nova instância do SubscriptionHandler
func NewSubscriptionHandler(service SubscriptionServiceInterface) *handler {
	return &handler{
		service: service,
	}
}

// writeErrorResponse escreve uma resposta de erro padronizada
func (h *handler) writeErrorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int, message string) {
	correlationID := logging.GetCorrelationID(r.Context())

	errorResponse := ErrorResponse{
		Error:         err.Error(),
		Message:       message,
		CorrelationID: correlationID,
		StatusCode:    statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Correlation-ID", correlationID)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// writeSuccessResponse escreve uma resposta de sucesso padronizada
func (h *handler) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int, message string) {
	correlationID := logging.GetCorrelationID(r.Context())

	successResponse := SuccessResponse{
		Data:          data,
		Message:       message,
		CorrelationID: correlationID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Correlation-ID", correlationID)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(successResponse)
}

// CreateSubscription handler para criar uma subscription
func (h *handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, r, err, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	subscription, err := h.service.CreateSubscription(r.Context(), req)
	if err != nil {
		h.writeErrorResponse(w, r, err, http.StatusInternalServerError, "Failed to create subscription")
		return
	}

	h.writeSuccessResponse(w, r, subscription, http.StatusCreated, "Subscription created successfully")
}

// GetSubscriptionByID handler para buscar uma subscription pelo ID
func (h *handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.writeErrorResponse(w, r,
			fmt.Errorf("missing subscription ID"),
			http.StatusBadRequest,
			"Subscription ID is required")
		return
	}

	subscription, err := h.service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		h.writeErrorResponse(w, r, err, http.StatusNotFound, "Subscription not found")
		return
	}

	h.writeSuccessResponse(w, r, subscription, http.StatusOK, "")
}

// GetAllSubscriptions handler para buscar todas as subscriptions
func (h *handler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.service.GetAllSubscriptions(r.Context())
	if err != nil {
		h.writeErrorResponse(w, r, err, http.StatusInternalServerError, "Failed to retrieve subscriptions")
		return
	}

	h.writeSuccessResponse(w, r, subscriptions, http.StatusOK, "")
}

// ActivateSubscription handler para ativar uma subscription
func (h *handler) ActivateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.writeErrorResponse(w, r,
			fmt.Errorf("missing subscription ID"),
			http.StatusBadRequest,
			"Subscription ID is required")
		return
	}

	// Tenta extrair correlation ID do header (opcional)
	correlationID := r.Header.Get("X-Correlation-ID")

	err := h.service.ActivateSubscription(r.Context(), id, correlationID)
	if err != nil {
		h.writeErrorResponse(w, r, err, http.StatusInternalServerError, "Failed to activate subscription")
		return
	}

	result := map[string]string{
		"message": "Subscription activated successfully",
		"id":      id,
	}
	h.writeSuccessResponse(w, r, result, http.StatusOK, "Subscription activated successfully")
}

// RegisterRoutes registra as rotas da subscription
func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscriptions", h.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions", h.GetAllSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", h.GetSubscriptionByID).Methods("GET")
	router.HandleFunc("/subscriptions/{id}/activate", h.ActivateSubscription).Methods("POST")
}
