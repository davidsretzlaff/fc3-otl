package subscription

import (
	"context"
	"fmt"
	"time"

	"payments-subscription/internal/common/logging"
	"payments-subscription/internal/customer"
)

// SubscriptionService representa o serviço de aplicação para Subscription
type SubscriptionService struct {
	repository     SubscriptionRepository
	eventService   *SubscriptionEventService
	customerClient *customer.CustomerClient
	logger         *logging.StructuredLogger
}

// NewSubscriptionService cria uma nova instância do SubscriptionService
func NewSubscriptionService(
	repository SubscriptionRepository,
	eventService *SubscriptionEventService,
	customerClient *customer.CustomerClient,
) *SubscriptionService {
	return &SubscriptionService{
		repository:     repository,
		eventService:   eventService,
		customerClient: customerClient,
		logger:         logging.NewStructuredLogger("subscription-service"),
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

// SubscriptionServiceInterface é uma interface para o SubscriptionService
type SubscriptionServiceInterface interface {
	CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*SubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, id string) (*SubscriptionResponse, error)
	GetAllSubscriptions(ctx context.Context) ([]*SubscriptionResponse, error)
	ActivateSubscription(ctx context.Context, id, correlationID string) error
}

// CreateSubscription cria uma nova subscription
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	startTime := time.Now()
	operation := "CreateSubscription"

	// Garantir que existe correlation ID
	ctx = logging.EnsureCorrelationID(ctx, "subscription")

	// Log início da operação
	s.logger.OperationStart(ctx, operation, map[string]interface{}{
		"plan_id":        req.PlanID,
		"customer_email": req.Customer.Email,
		"customer_name":  req.Customer.Name,
	})

	// Criar o customer primeiro
	customerReq := customer.CustomerRequest{
		Name:  req.Customer.Name,
		Email: req.Customer.Email,
	}

	s.logger.Info(ctx, operation, "Iniciando criação de customer", map[string]interface{}{
		"customer_email": req.Customer.Email,
		"customer_name":  req.Customer.Name,
	})

	customerResp, err := s.customerClient.CreateCustomer(ctx, customerReq)
	if err != nil {
		s.logger.Error(ctx, operation, "Erro ao criar customer", err, map[string]interface{}{
			"customer_email": req.Customer.Email,
			"plan_id":        req.PlanID,
		})
		return nil, fmt.Errorf("erro ao criar customer: %w", err)
	}

	s.logger.Info(ctx, operation, "Customer criado com sucesso", map[string]interface{}{
		"customer_id":    customerResp.ID,
		"customer_email": customerResp.Email,
	})

	// Criar a subscription usando o ID do customer retornado
	correlationID := logging.GetCorrelationID(ctx)
	subscription, err := NewSubscription(req.PlanID, customerResp.ID, correlationID)
	if err != nil {
		s.logger.Error(ctx, operation, "Erro ao criar entidade subscription", err, map[string]interface{}{
			"plan_id":     req.PlanID,
			"customer_id": customerResp.ID,
		})
		return nil, fmt.Errorf("erro ao criar subscription: %w", err)
	}

	s.logger.Info(ctx, operation, "Salvando subscription no repositório", map[string]interface{}{
		"subscription_id": subscription.ID().String(),
		"customer_id":     customerResp.ID,
		"plan_id":         req.PlanID,
		"status":          string(subscription.Status()),
	})

	if err := s.repository.Create(ctx, subscription); err != nil {
		s.logger.Error(ctx, operation, "Erro ao salvar subscription no banco", err, map[string]interface{}{
			"subscription_id": subscription.ID().String(),
			"customer_id":     customerResp.ID,
			"plan_id":         req.PlanID,
		})
		return nil, fmt.Errorf("erro ao salvar subscription: %w", err)
	}

	/*if err := s.eventService.PublishSubscriptionEvents(ctx, subscription); err != nil {
		s.logger.Warn(ctx, operation, "Erro ao publicar eventos", map[string]interface{}{
			"subscription_id": subscription.ID().String(),
			"error":          err.Error(),
		})
	}*/

	response := s.toSubscriptionResponse(subscription)

	// Log fim da operação
	s.logger.OperationEnd(ctx, operation, startTime, map[string]interface{}{
		"subscription_id": response.ID,
		"customer_id":     response.CustomerID,
		"plan_id":         response.PlanID,
		"status":          response.Status,
	})

	return response, nil
}

// GetSubscriptionByID busca uma subscription pelo ID
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id string) (*SubscriptionResponse, error) {
	startTime := time.Now()
	operation := "GetSubscriptionByID"

	ctx = logging.EnsureCorrelationID(ctx, "subscription")

	s.logger.OperationStart(ctx, operation, map[string]interface{}{
		"subscription_id": id,
	})

	subscriptionID, err := NewSubscriptionIDFromString(id)
	if err != nil {
		s.logger.Error(ctx, operation, "ID de subscription inválido", err, map[string]interface{}{
			"provided_id": id,
		})
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	subscription, err := s.repository.GetByID(ctx, subscriptionID)
	if err != nil {
		s.logger.Error(ctx, operation, "Erro ao buscar subscription no banco", err, map[string]interface{}{
			"subscription_id": id,
		})
		return nil, fmt.Errorf("erro ao buscar subscription: %w", err)
	}

	response := s.toSubscriptionResponse(subscription)

	s.logger.OperationEnd(ctx, operation, startTime, map[string]interface{}{
		"subscription_id": response.ID,
		"customer_id":     response.CustomerID,
		"status":          response.Status,
	})

	return response, nil
}

// GetAllSubscriptions busca todas as subscriptions
func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]*SubscriptionResponse, error) {
	startTime := time.Now()
	operation := "GetAllSubscriptions"

	ctx = logging.EnsureCorrelationID(ctx, "subscription")

	s.logger.OperationStart(ctx, operation, nil)

	subscriptions, err := s.repository.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, operation, "Erro ao buscar todas as subscriptions", err, nil)
		return nil, fmt.Errorf("erro ao buscar subscriptions: %w", err)
	}

	responses := make([]*SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = s.toSubscriptionResponse(subscription)
	}

	s.logger.OperationEnd(ctx, operation, startTime, map[string]interface{}{
		"total_found": len(responses),
	})

	return responses, nil
}

// ActivateSubscription ativa uma subscription
func (s *SubscriptionService) ActivateSubscription(ctx context.Context, id, correlationID string) error {
	startTime := time.Now()
	operation := "ActivateSubscription"

	// Usar o correlation ID fornecido ou gerar um novo
	if correlationID != "" {
		ctx = logging.WithCorrelationID(ctx, correlationID)
	} else {
		ctx = logging.EnsureCorrelationID(ctx, "subscription")
	}

	s.logger.OperationStart(ctx, operation, map[string]interface{}{
		"subscription_id": id,
		"correlation_id":  correlationID,
	})

	subscriptionID, err := NewSubscriptionIDFromString(id)
	if err != nil {
		s.logger.Error(ctx, operation, "ID de subscription inválido", err, map[string]interface{}{
			"provided_id": id,
		})
		return fmt.Errorf("ID inválido: %w", err)
	}

	subscription, err := s.repository.GetByID(ctx, subscriptionID)
	if err != nil {
		s.logger.Error(ctx, operation, "Erro ao buscar subscription", err, map[string]interface{}{
			"subscription_id": id,
		})
		return fmt.Errorf("erro ao buscar subscription: %w", err)
	}

	currentCorrelationID := logging.GetCorrelationID(ctx)
	if err := subscription.Activate(currentCorrelationID); err != nil {
		s.logger.Error(ctx, operation, "Erro na ativação da subscription", err, map[string]interface{}{
			"subscription_id": id,
			"current_status":  string(subscription.Status()),
		})
		return fmt.Errorf("erro ao ativar subscription: %w", err)
	}

	if err := s.repository.Update(ctx, subscription); err != nil {
		s.logger.Error(ctx, operation, "Erro ao atualizar subscription no banco", err, map[string]interface{}{
			"subscription_id": id,
			"new_status":      string(subscription.Status()),
		})
		return fmt.Errorf("erro ao atualizar subscription: %w", err)
	}

	s.logger.OperationEnd(ctx, operation, startTime, map[string]interface{}{
		"subscription_id": id,
		"new_status":      string(subscription.Status()),
	})

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
