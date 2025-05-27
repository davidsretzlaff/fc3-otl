package subscription

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Subscription representa o agregado principal do domínio de Subscription
type Subscription struct {
	id         SubscriptionID
	planID     PlanID
	customerID CustomerID
	status     SubscriptionStatus
	createdAt  time.Time
	updatedAt  time.Time
	events     []DomainEvent
}

// SubscriptionID é um value object para o ID da subscription
type SubscriptionID struct {
	value string
}

// PlanID é um value object para o ID do plano
type PlanID struct {
	value string
}

// CustomerID é um value object para o ID do customer
type CustomerID struct {
	value string
}

// SubscriptionStatus representa o status da subscription
type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusInactive  SubscriptionStatus = "inactive"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
)

// Erros do domínio
var (
	ErrInvalidPlanID           = errors.New("plan ID é obrigatório")
	ErrInvalidCustomerID       = errors.New("customer ID é obrigatório")
	ErrSubscriptionNotFound    = errors.New("subscription não encontrada")
	ErrInvalidStatusTransition = errors.New("transição de status inválida")
)

// DomainEvent representa um evento de domínio
type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	CorrelationID() string
}

// BaseEvent implementa campos comuns dos eventos
type BaseEvent struct {
	eventType     string
	aggregateID   string
	occurredAt    time.Time
	correlationID string
}

func (e BaseEvent) EventType() string     { return e.eventType }
func (e BaseEvent) AggregateID() string   { return e.aggregateID }
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) CorrelationID() string { return e.correlationID }

// Eventos de domínio
type SubscriptionRequestedEvent struct {
	BaseEvent
	PlanID     string `json:"plan_id"`
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
}

type SubscriptionReadyForActivationEvent struct {
	BaseEvent
	PlanID     string `json:"plan_id"`
	CustomerID string `json:"customer_id"`
}

type SubscriptionActivatedEvent struct {
	BaseEvent
	PlanID     string `json:"plan_id"`
	CustomerID string `json:"customer_id"`
}

type SubscriptionCancelledEvent struct {
	BaseEvent
	Reason string `json:"reason"`
}

type SubscriptionSuspendedEvent struct {
	BaseEvent
	Reason string `json:"reason"`
}

// NewSubscriptionID cria um novo SubscriptionID
func NewSubscriptionID() SubscriptionID {
	return SubscriptionID{value: uuid.New().String()}
}

// NewSubscriptionIDFromString cria um SubscriptionID a partir de uma string
func NewSubscriptionIDFromString(id string) (SubscriptionID, error) {
	if id == "" {
		return SubscriptionID{}, errors.New("ID não pode ser vazio")
	}
	return SubscriptionID{value: id}, nil
}

// String retorna a representação em string do SubscriptionID
func (s SubscriptionID) String() string {
	return s.value
}

// NewPlanID cria um PlanID a partir de uma string
func NewPlanID(id string) (PlanID, error) {
	if id == "" {
		return PlanID{}, ErrInvalidPlanID
	}
	return PlanID{value: id}, nil
}

// String retorna a representação em string do PlanID
func (p PlanID) String() string {
	return p.value
}

// NewCustomerID cria um CustomerID
func NewCustomerID(id string) (CustomerID, error) {
	if id == "" {
		return CustomerID{value: uuid.New().String()}, nil
	}
	return CustomerID{value: id}, nil
}

// String retorna a representação em string do CustomerID
func (c CustomerID) String() string {
	return c.value
}

// NewSubscription cria uma nova subscription
func NewSubscription(planID, customerID string, correlationID string) (*Subscription, error) {
	if err := validateSubscriptionData(planID); err != nil {
		return nil, err
	}

	pID, _ := NewPlanID(planID)
	cID, _ := NewCustomerID(customerID)

	subscription := &Subscription{
		id:         NewSubscriptionID(),
		planID:     pID,
		customerID: cID,
		status:     SubscriptionStatusPending,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
		events:     make([]DomainEvent, 0),
	}

	// Adiciona evento de subscription solicitada
	event := SubscriptionRequestedEvent{
		BaseEvent: BaseEvent{
			eventType:     "SubscriptionRequested",
			aggregateID:   subscription.id.String(),
			occurredAt:    time.Now(),
			correlationID: correlationID,
		},
		PlanID:     planID,
		CustomerID: customerID,
	}

	subscription.addEvent(event)
	return subscription, nil
}

// ReconstructSubscription reconstrói uma subscription a partir de dados persistidos
func ReconstructSubscription(id, planID, customerID string, status SubscriptionStatus, createdAt, updatedAt time.Time) (*Subscription, error) {
	subscriptionID, err := NewSubscriptionIDFromString(id)
	if err != nil {
		return nil, err
	}

	pID, err := NewPlanID(planID)
	if err != nil {
		return nil, err
	}

	cID, err := NewCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	return &Subscription{
		id:         subscriptionID,
		planID:     pID,
		customerID: cID,
		status:     status,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		events:     make([]DomainEvent, 0),
	}, nil
}

// Getters
func (s *Subscription) ID() SubscriptionID {
	return s.id
}

func (s *Subscription) PlanID() PlanID {
	return s.planID
}

func (s *Subscription) CustomerID() CustomerID {
	return s.customerID
}

func (s *Subscription) Status() SubscriptionStatus {
	return s.status
}

func (s *Subscription) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Subscription) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Subscription) Events() []DomainEvent {
	return s.events
}

// MarkAsReadyForActivation marca a subscription como pronta para ativação
func (s *Subscription) MarkAsReadyForActivation(correlationID string) error {
	if s.status != SubscriptionStatusPending {
		return ErrInvalidStatusTransition
	}

	event := SubscriptionReadyForActivationEvent{
		BaseEvent: BaseEvent{
			eventType:     "SubscriptionReadyForActivation",
			aggregateID:   s.id.String(),
			occurredAt:    time.Now(),
			correlationID: correlationID,
		},
		PlanID:     s.planID.String(),
		CustomerID: s.customerID.String(),
	}

	s.addEvent(event)
	s.updatedAt = time.Now()
	return nil
}

// Activate ativa a subscription
func (s *Subscription) Activate(correlationID string) error {
	if s.status != SubscriptionStatusPending && s.status != SubscriptionStatusSuspended {
		return ErrInvalidStatusTransition
	}

	s.status = SubscriptionStatusActive
	s.updatedAt = time.Now()

	event := SubscriptionActivatedEvent{
		BaseEvent: BaseEvent{
			eventType:     "SubscriptionActivated",
			aggregateID:   s.id.String(),
			occurredAt:    time.Now(),
			correlationID: correlationID,
		},
		PlanID:     s.planID.String(),
		CustomerID: s.customerID.String(),
	}

	s.addEvent(event)
	return nil
}

// Cancel cancela a subscription
func (s *Subscription) Cancel(reason, correlationID string) error {
	if s.status == SubscriptionStatusCancelled {
		return ErrInvalidStatusTransition
	}

	s.status = SubscriptionStatusCancelled
	s.updatedAt = time.Now()

	event := SubscriptionCancelledEvent{
		BaseEvent: BaseEvent{
			eventType:     "SubscriptionCancelled",
			aggregateID:   s.id.String(),
			occurredAt:    time.Now(),
			correlationID: correlationID,
		},
		Reason: reason,
	}

	s.addEvent(event)
	return nil
}

// Suspend suspende a subscription
func (s *Subscription) Suspend(reason, correlationID string) error {
	if s.status != SubscriptionStatusActive {
		return ErrInvalidStatusTransition
	}

	s.status = SubscriptionStatusSuspended
	s.updatedAt = time.Now()

	event := SubscriptionSuspendedEvent{
		BaseEvent: BaseEvent{
			eventType:     "SubscriptionSuspended",
			aggregateID:   s.id.String(),
			occurredAt:    time.Now(),
			correlationID: correlationID,
		},
		Reason: reason,
	}

	s.addEvent(event)
	return nil
}

// IsActive verifica se a subscription está ativa
func (s *Subscription) IsActive() bool {
	return s.status == SubscriptionStatusActive
}

// IsPending verifica se a subscription está pendente
func (s *Subscription) IsPending() bool {
	return s.status == SubscriptionStatusPending
}

// ClearEvents limpa os eventos (usado após persistência)
func (s *Subscription) ClearEvents() {
	s.events = make([]DomainEvent, 0)
}

// addEvent adiciona um evento à lista de eventos
func (s *Subscription) addEvent(event DomainEvent) {
	s.events = append(s.events, event)
}

// validateSubscriptionData valida os dados da subscription
func validateSubscriptionData(planID string) error {
	if planID == "" {
		return ErrInvalidPlanID
	}
	return nil
}

// SubscriptionRepository define o contrato para persistência de subscriptions
type SubscriptionRepository interface {
	// Create cria uma nova subscription
	Create(ctx context.Context, subscription *Subscription) error

	// GetByID busca uma subscription pelo ID
	GetByID(ctx context.Context, id SubscriptionID) (*Subscription, error)

	// GetByCustomerID busca subscriptions pelo customer ID
	GetByCustomerID(ctx context.Context, customerID CustomerID) ([]*Subscription, error)

	// Update atualiza uma subscription existente
	Update(ctx context.Context, subscription *Subscription) error

	// GetAll busca todas as subscriptions
	GetAll(ctx context.Context) ([]*Subscription, error)
}
