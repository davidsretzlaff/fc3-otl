package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payments-customer/internal/customer"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// MySQLCustomerRepository implementa o CustomerRepository usando MySQL
type MySQLCustomerRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

// NewMySQLCustomerRepository cria uma nova instância do repositório MySQL
func NewMySQLCustomerRepository(db *sql.DB, tracer trace.Tracer) *MySQLCustomerRepository {
	return &MySQLCustomerRepository{
		db:     db,
		tracer: tracer,
	}
}

// Create cria um novo customer no banco de dados
func (r *MySQLCustomerRepository) Create(ctx context.Context, customer *customer.Customer) error {
	ctx, span := r.tracer.Start(ctx, "MySQLCustomerRepository.Create")
	defer span.End()

	query := `
		INSERT INTO customers (id, name, email, document, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		customer.ID().String(),
		customer.Name(),
		customer.Email(),
		customer.Document(),
		string(customer.Status()),
		customer.CreatedAt(),
		customer.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir customer no banco: %w", err)
	}

	return nil
}

// GetByID busca um customer pelo ID no banco de dados
func (r *MySQLCustomerRepository) GetByID(ctx context.Context, id customer.CustomerID) (*customer.Customer, error) {
	ctx, span := r.tracer.Start(ctx, "MySQLCustomerRepository.GetByID")
	defer span.End()

	query := `
		SELECT id, name, email, document, status, created_at, updated_at
		FROM customers
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id.String())

	var customerID, name, email, document, status string
	var createdAt, updatedAt time.Time

	err := row.Scan(&customerID, &name, &email, &document, &status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("erro ao buscar customer no banco: %w", err)
	}

	customerEntity, err := customer.ReconstructCustomer(
		customerID,
		name,
		email,
		document,
		customer.CustomerStatus(status),
		createdAt,
		updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao reconstruir customer: %w", err)
	}

	return customerEntity, nil
}
