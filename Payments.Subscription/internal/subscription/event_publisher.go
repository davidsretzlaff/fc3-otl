package subscription

import (
	"context"
	"encoding/json"
	"fmt"

	"payments-subscription/internal/common/logging"
)

// EventPublisher define o contrato para publicação de eventos
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// InMemoryEventPublisher implementa um publisher simples em memória para demonstração
type InMemoryEventPublisher struct {
	logger *logging.StructuredLogger
}

// NewInMemoryEventPublisher cria uma nova instância do publisher
func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		logger: logging.NewStructuredLogger("subscription-service"),
	}
}

// Publish publica um evento (implementação simples para demonstração)
func (p *InMemoryEventPublisher) Publish(ctx context.Context, event DomainEvent) error {
	// Serializa o evento para JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	// Log do evento usando logger estruturado
	p.logger.Info(ctx, "EventPublished",
		fmt.Sprintf("Event published: %s", event.EventType()),
		map[string]interface{}{
			"event_type": event.EventType(),
			"event_data": string(eventData),
		})

	return nil
}

// EventHandler define o contrato para manipuladores de eventos
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	CanHandle(eventType string) bool
}

// SubscriptionEventService gerencia a publicação de eventos de subscription
type SubscriptionEventService struct {
	publisher EventPublisher
}

// NewSubscriptionEventService cria uma nova instância do serviço de eventos
func NewSubscriptionEventService(publisher EventPublisher) *SubscriptionEventService {
	return &SubscriptionEventService{
		publisher: publisher,
	}
}

// PublishSubscriptionEvents publica todos os eventos pendentes de uma subscription
func (s *SubscriptionEventService) PublishSubscriptionEvents(ctx context.Context, subscription *Subscription) error {

	events := subscription.Events()
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		if err := s.publisher.Publish(ctx, event); err != nil {
			return fmt.Errorf("erro ao publicar evento %s: %w", event.EventType(), err)
		}
	}

	// Limpa os eventos após publicação
	subscription.ClearEvents()
	return nil
}
