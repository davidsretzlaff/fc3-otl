package customer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

// CustomerService representa o serviço de aplicação para Customer
type CustomerService struct {
	repository CustomerRepository
	tracer     trace.Tracer
}

// NewCustomerService cria uma nova instância do CustomerService
func NewCustomerService(repository CustomerRepository, tracer trace.Tracer) *CustomerService {
	return &CustomerService{
		repository: repository,
		tracer:     tracer,
	}
}

// CreateCustomerRequest representa a requisição para criar um customer
type CreateCustomerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
}

// CustomerResponse representa a resposta com dados do customer
type CustomerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Document  string `json:"document"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateCustomer cria um novo customer
func (s *CustomerService) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*CustomerResponse, error) {
	ctx, span := s.tracer.Start(ctx, "CustomerService.CreateCustomer")
	defer span.End()

	// Cria o customer
	customer, err := NewCustomer(req.Name, req.Email, req.Document)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar customer: %w", err)
	}

	// Salva no repositório
	if err := s.repository.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("erro ao salvar customer: %w", err)
	}

	return s.toCustomerResponse(customer), nil
}

// GetCustomerByID busca um customer pelo ID
func (s *CustomerService) GetCustomerByID(ctx context.Context, id string) (*CustomerResponse, error) {
	ctx, span := s.tracer.Start(ctx, "CustomerService.GetCustomerByID")
	defer span.End()

	customerID, err := NewCustomerIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	customer, err := s.repository.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar customer: %w", err)
	}

	return s.toCustomerResponse(customer), nil
}

// toCustomerResponse converte um Customer para CustomerResponse
func (s *CustomerService) toCustomerResponse(customer *Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:        customer.ID().String(),
		Name:      customer.Name(),
		Email:     customer.Email(),
		Document:  customer.Document(),
		Status:    string(customer.Status()),
		CreatedAt: customer.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: customer.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
