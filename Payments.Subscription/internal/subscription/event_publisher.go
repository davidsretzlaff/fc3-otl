package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// EventPublisher define o contrato para publicação de eventos
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// InMemoryEventPublisher implementa um publisher simples em memória para demonstração
type InMemoryEventPublisher struct {
}

// NewInMemoryEventPublisher cria uma nova instância do publisher
func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{}
}

// Publish publica um evento (implementação simples para demonstração)
func (p *InMemoryEventPublisher) Publish(ctx context.Context, event DomainEvent) error {
	// Serializa o evento para JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	// Log do evento (em um cenário real, seria enviado para um message broker)
	log.Printf("Evento publicado: %s - %s - CorrelationID: %s",
		event.EventType(),
		string(eventData),
		event.CorrelationID())

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
