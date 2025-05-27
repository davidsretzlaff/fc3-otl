package subscription

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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

// CreateSubscription handler para criar uma subscription
func (h *handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.CreateSubscription(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GetSubscriptionByID handler para buscar uma subscription pelo ID
func (h *handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// GetAllSubscriptions handler para buscar todas as subscriptions
func (h *handler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {

	subscriptions, err := h.service.GetAllSubscriptions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// RegisterRoutes registra as rotas da subscription
func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscriptions", h.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions", h.GetAllSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", h.GetSubscriptionByID).Methods("GET")
}
