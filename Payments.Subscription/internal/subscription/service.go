package subscription

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/trace"
)

// SubscriptionService representa o serviço de aplicação para Subscription
type SubscriptionService struct {
	repository   SubscriptionRepository
	eventService *SubscriptionEventService
	tracer       trace.Tracer
}

// NewSubscriptionService cria uma nova instância do SubscriptionService
func NewSubscriptionService(repository SubscriptionRepository, eventService *SubscriptionEventService, tracer trace.Tracer) *SubscriptionService {
	return &SubscriptionService{
		repository:   repository,
		eventService: eventService,
		tracer:       tracer,
	}
}

// CreateSubscriptionRequest representa a requisição para criar uma subscription
type CreateSubscriptionRequest struct {
	PlanID   string                `json:"plan_id"`
	Customer CreateCustomerRequest `json:"customer"`
}

// CreateCustomerRequest representa a requisição para criar um customer
type CreateCustomerRequest struct {
	CustomerID string `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

// SubscriptionResponse representa a resposta com dados da subscription
type SubscriptionResponse struct {
	ID         string `json:"id"`
	PlanID     string `json:"plan_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// CreateSubscription cria uma nova subscription
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "SubscriptionService.CreateSubscription")
	defer span.End()

	// Define correlation ID se não fornecido
	correlationID := fmt.Sprintf("sub-%d", ctx.Value("request_id"))

	// Cria a subscription
	subscription, err := NewSubscription(req.PlanID, req.Customer.CustomerID, correlationID)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar subscription: %w", err)
	}

	// Salva no repositório
	if err := s.repository.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("erro ao salvar subscription: %w", err)
	}

	// Publica os eventos de domínio
	if err := s.eventService.PublishSubscriptionEvents(ctx, subscription); err != nil {
		// Log do erro mas não falha a operação
		log.Printf("Erro ao publicar eventos: %v", err)
	}

	return s.toSubscriptionResponse(subscription), nil
}

// GetSubscriptionByID busca uma subscription pelo ID
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id string) (*SubscriptionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "SubscriptionService.GetSubscriptionByID")
	defer span.End()

	subscriptionID, err := NewSubscriptionIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	subscription, err := s.repository.GetByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar subscription: %w", err)
	}

	return s.toSubscriptionResponse(subscription), nil
}

// GetAllSubscriptions busca todas as subscriptions
func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]*SubscriptionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "SubscriptionService.GetAllSubscriptions")
	defer span.End()

	subscriptions, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar subscriptions: %w", err)
	}

	responses := make([]*SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = s.toSubscriptionResponse(subscription)
	}

	return responses, nil
}

// ActivateSubscription ativa uma subscription
func (s *SubscriptionService) ActivateSubscription(ctx context.Context, id, correlationID string) error {
	ctx, span := s.tracer.Start(ctx, "SubscriptionService.ActivateSubscription")
	defer span.End()

	subscriptionID, err := NewSubscriptionIDFromString(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}

	subscription, err := s.repository.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("erro ao buscar subscription: %w", err)
	}

	if err := subscription.Activate(correlationID); err != nil {
		return fmt.Errorf("erro ao ativar subscription: %w", err)
	}

	if err := s.repository.Update(ctx, subscription); err != nil {
		return fmt.Errorf("erro ao atualizar subscription: %w", err)
	}

	return nil
}

// toSubscriptionResponse converte uma Subscription para SubscriptionResponse
func (s *SubscriptionService) toSubscriptionResponse(subscription *Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:         subscription.ID().String(),
		PlanID:     subscription.PlanID().String(),
		CustomerID: subscription.CustomerID().String(),
		Status:     string(subscription.Status()),
		CreatedAt:  subscription.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  subscription.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
