package customer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
)

// handler gerencia as requisições HTTP para Customer
type handler struct {
	service *CustomerService
	tracer  trace.Tracer
}

// NewCustomerHandler cria uma nova instância do CustomerHandler
func NewCustomerHandler(service *CustomerService, tracer trace.Tracer) *handler {
	return &handler{
		service: service,
		tracer:  tracer,
	}
}

// CreateCustomer handler para criar um customer
func (h *handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "CustomerHandler.CreateCustomer")
	defer span.End()

	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	customer, err := h.service.CreateCustomer(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetCustomerByID handler para buscar um customer pelo ID
func (h *handler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "CustomerHandler.GetCustomerByID")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	customer, err := h.service.GetCustomerByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// RegisterRoutes registra as rotas do customer
func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/customers", h.CreateCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", h.GetCustomerByID).Methods("GET")
}
