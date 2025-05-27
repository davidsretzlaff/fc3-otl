package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// EventPublisher define o contrato para publicação de eventos
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// InMemoryEventPublisher implementa um publisher simples em memória para demonstração
type InMemoryEventPublisher struct {
	tracer trace.Tracer
}

// NewInMemoryEventPublisher cria uma nova instância do publisher
func NewInMemoryEventPublisher(tracer trace.Tracer) *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		tracer: tracer,
	}
}

// Publish publica um evento (implementação simples para demonstração)
func (p *InMemoryEventPublisher) Publish(ctx context.Context, event DomainEvent) error {
	ctx, span := p.tracer.Start(ctx, "EventPublisher.Publish")
	defer span.End()

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

	// Adiciona informações ao span
	span.SetAttributes(
		attribute.String("event.type", event.EventType()),
		attribute.String("event.aggregate_id", event.AggregateID()),
		attribute.String("event.correlation_id", event.CorrelationID()),
	)

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
	tracer    trace.Tracer
}

// NewSubscriptionEventService cria uma nova instância do serviço de eventos
func NewSubscriptionEventService(publisher EventPublisher, tracer trace.Tracer) *SubscriptionEventService {
	return &SubscriptionEventService{
		publisher: publisher,
		tracer:    tracer,
	}
}

// PublishSubscriptionEvents publica todos os eventos pendentes de uma subscription
func (s *SubscriptionEventService) PublishSubscriptionEvents(ctx context.Context, subscription *Subscription) error {
	ctx, span := s.tracer.Start(ctx, "SubscriptionEventService.PublishSubscriptionEvents")
	defer span.End()

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
