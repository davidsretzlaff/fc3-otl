package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payments-customer/internal/subscription"
	"time"
)

// MySQLSubscriptionRepository implementa o SubscriptionRepository usando MySQL
type MySQLSubscriptionRepository struct {
	db *sql.DB
}

// NewMySQLSubscriptionRepository cria uma nova instância do repositório MySQL
func NewMySQLSubscriptionRepository(db *sql.DB) *MySQLSubscriptionRepository {
	return &MySQLSubscriptionRepository{
		db: db,
	}
}

// Create cria uma nova subscription no banco de dados
func (r *MySQLSubscriptionRepository) Create(ctx context.Context, sub *subscription.Subscription) error {

	query := `
		INSERT INTO subscriptions (id, plan_id, customer_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		sub.ID().String(),
		sub.PlanID().String(),
		sub.CustomerID().String(),
		string(sub.Status()),
		sub.CreatedAt(),
		sub.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir subscription no banco: %w", err)
	}

	return nil
}

// GetByID busca uma subscription pelo ID no banco de dados
func (r *MySQLSubscriptionRepository) GetByID(ctx context.Context, id subscription.SubscriptionID) (*subscription.Subscription, error) {
	query := `
		SELECT id, plan_id, customer_id, status, created_at, updated_at
		FROM subscriptions
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id.String())

	var subscriptionID, planID, customerID, status string
	var createdAt, updatedAt time.Time

	err := row.Scan(&subscriptionID, &planID, &customerID, &status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, subscription.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("erro ao buscar subscription no banco: %w", err)
	}

	subscriptionEntity, err := subscription.ReconstructSubscription(
		subscriptionID,
		planID,
		customerID,
		subscription.SubscriptionStatus(status),
		createdAt,
		updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao reconstruir subscription: %w", err)
	}

	return subscriptionEntity, nil
}

// GetByCustomerID busca subscriptions pelo customer ID no banco de dados
func (r *MySQLSubscriptionRepository) GetByCustomerID(ctx context.Context, customerID subscription.CustomerID) ([]*subscription.Subscription, error) {

	query := `
		SELECT id, plan_id, customer_id, status, created_at, updated_at
		FROM subscriptions
		WHERE customer_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, customerID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar subscriptions no banco: %w", err)
	}
	defer rows.Close()

	var subscriptions []*subscription.Subscription

	for rows.Next() {
		var subscriptionID, planID, custID, status string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&subscriptionID, &planID, &custID, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer scan da subscription: %w", err)
		}

		subscriptionEntity, err := subscription.ReconstructSubscription(
			subscriptionID,
			planID,
			custID,
			subscription.SubscriptionStatus(status),
			createdAt,
			updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao reconstruir subscription: %w", err)
		}

		subscriptions = append(subscriptions, subscriptionEntity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre as subscriptions: %w", err)
	}

	return subscriptions, nil
}

// Update atualiza uma subscription existente no banco de dados
func (r *MySQLSubscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {

	query := `
		UPDATE subscriptions 
		SET plan_id = ?, customer_id = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		sub.PlanID().String(),
		sub.CustomerID().String(),
		string(sub.Status()),
		sub.UpdatedAt(),
		sub.ID().String(),
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar subscription no banco: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return subscription.ErrSubscriptionNotFound
	}

	return nil
}

// GetAll busca todas as subscriptions no banco de dados
func (r *MySQLSubscriptionRepository) GetAll(ctx context.Context) ([]*subscription.Subscription, error) {

	query := `
		SELECT id, plan_id, customer_id, status, created_at, updated_at
		FROM subscriptions
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar subscriptions no banco: %w", err)
	}
	defer rows.Close()

	var subscriptions []*subscription.Subscription

	for rows.Next() {
		var subscriptionID, planID, customerID, status string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&subscriptionID, &planID, &customerID, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer scan da subscription: %w", err)
		}

		subscriptionEntity, err := subscription.ReconstructSubscription(
			subscriptionID,
			planID,
			customerID,
			subscription.SubscriptionStatus(status),
			createdAt,
			updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao reconstruir subscription: %w", err)
		}

		subscriptions = append(subscriptions, subscriptionEntity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre as subscriptions: %w", err)
	}

	return subscriptions, nil
}
