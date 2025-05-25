package customer

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Customer representa o agregado principal do domínio de Customer
type Customer struct {
	id        CustomerID
	name      string
	email     string
	document  string
	status    CustomerStatus
	createdAt time.Time
	updatedAt time.Time
}

// CustomerID é um value object para o ID do customer
type CustomerID struct {
	value string
}

// CustomerStatus representa o status do customer
type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
	CustomerStatusBlocked  CustomerStatus = "blocked"
)

// Erros do domínio
var (
	ErrInvalidCustomerName     = errors.New("nome do customer é obrigatório")
	ErrInvalidCustomerEmail    = errors.New("email do customer é inválido")
	ErrInvalidCustomerDocument = errors.New("documento do customer é inválido")
	ErrCustomerNotFound        = errors.New("customer não encontrado")
)

// NewCustomerID cria um novo CustomerID
func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New().String()}
}

// NewCustomerIDFromString cria um CustomerID a partir de uma string
func NewCustomerIDFromString(id string) (CustomerID, error) {
	if id == "" {
		return CustomerID{}, errors.New("ID não pode ser vazio")
	}
	return CustomerID{value: id}, nil
}

// String retorna a representação em string do CustomerID
func (c CustomerID) String() string {
	return c.value
}

// NewCustomer cria um novo customer
func NewCustomer(name, email, document string) (*Customer, error) {
	if err := validateCustomerData(name, email, document); err != nil {
		return nil, err
	}

	return &Customer{
		id:        NewCustomerID(),
		name:      name,
		email:     email,
		document:  document,
		status:    CustomerStatusActive,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ReconstructCustomer reconstrói um customer a partir de dados persistidos
func ReconstructCustomer(id, name, email, document string, status CustomerStatus, createdAt, updatedAt time.Time) (*Customer, error) {
	customerID, err := NewCustomerIDFromString(id)
	if err != nil {
		return nil, err
	}

	return &Customer{
		id:        customerID,
		name:      name,
		email:     email,
		document:  document,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// Getters
func (c *Customer) ID() CustomerID {
	return c.id
}

func (c *Customer) Name() string {
	return c.name
}

func (c *Customer) Email() string {
	return c.email
}

func (c *Customer) Document() string {
	return c.document
}

func (c *Customer) Status() CustomerStatus {
	return c.status
}

func (c *Customer) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Customer) UpdatedAt() time.Time {
	return c.updatedAt
}

// UpdateName atualiza o nome do customer
func (c *Customer) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidCustomerName
	}
	c.name = name
	c.updatedAt = time.Now()
	return nil
}

// UpdateEmail atualiza o email do customer
func (c *Customer) UpdateEmail(email string) error {
	if !isValidEmail(email) {
		return ErrInvalidCustomerEmail
	}
	c.email = email
	c.updatedAt = time.Now()
	return nil
}

// Block bloqueia o customer
func (c *Customer) Block() {
	c.status = CustomerStatusBlocked
	c.updatedAt = time.Now()
}

// Activate ativa o customer
func (c *Customer) Activate() {
	c.status = CustomerStatusActive
	c.updatedAt = time.Now()
}

// Deactivate desativa o customer
func (c *Customer) Deactivate() {
	c.status = CustomerStatusInactive
	c.updatedAt = time.Now()
}

// IsActive verifica se o customer está ativo
func (c *Customer) IsActive() bool {
	return c.status == CustomerStatusActive
}

// validateCustomerData valida os dados do customer
func validateCustomerData(name, email, document string) error {
	if name == "" {
		return ErrInvalidCustomerName
	}
	if !isValidEmail(email) {
		return ErrInvalidCustomerEmail
	}
	if document == "" {
		return ErrInvalidCustomerDocument
	}
	return nil
}

// isValidEmail valida se o email é válido (implementação simples)
func isValidEmail(email string) bool {
	return email != "" && len(email) > 3 && len(email) < 255
}

// CustomerRepository define o contrato para persistência de customers
type CustomerRepository interface {
	// Create cria um novo customer
	Create(ctx context.Context, customer *Customer) error

	// GetByID busca um customer pelo ID
	GetByID(ctx context.Context, id CustomerID) (*Customer, error)
}
