package subscription

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// SubscriptionServiceTracingDecorator é um decorator que adiciona tracing ao serviço
type SubscriptionServiceTracingDecorator struct {
	service SubscriptionServiceInterface
	tracer  trace.Tracer
}

// NewSubscriptionServiceTracingDecorator cria uma nova instância do decorator de tracing
func NewSubscriptionServiceTracingDecorator(service SubscriptionServiceInterface, tracer trace.Tracer) SubscriptionServiceInterface {
	return &SubscriptionServiceTracingDecorator{
		service: service,
		tracer:  tracer,
	}
}

// CreateSubscription adiciona tracing à operação de criação
func (d *SubscriptionServiceTracingDecorator) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	ctx, span := d.tracer.Start(ctx, "SubscriptionService.CreateSubscription")
	defer span.End()

	return d.service.CreateSubscription(ctx, req)
}

// GetSubscriptionByID adiciona tracing à operação de busca por ID
func (d *SubscriptionServiceTracingDecorator) GetSubscriptionByID(ctx context.Context, id string) (*SubscriptionResponse, error) {
	ctx, span := d.tracer.Start(ctx, "SubscriptionService.GetSubscriptionByID")
	defer span.End()

	return d.service.GetSubscriptionByID(ctx, id)
}

// GetAllSubscriptions adiciona tracing à operação de busca de todas as subscriptions
func (d *SubscriptionServiceTracingDecorator) GetAllSubscriptions(ctx context.Context) ([]*SubscriptionResponse, error) {
	ctx, span := d.tracer.Start(ctx, "SubscriptionService.GetAllSubscriptions")
	defer span.End()

	return d.service.GetAllSubscriptions(ctx)
}

// ActivateSubscription adiciona tracing à operação de ativação
func (d *SubscriptionServiceTracingDecorator) ActivateSubscription(ctx context.Context, id, correlationID string) error {
	ctx, span := d.tracer.Start(ctx, "SubscriptionService.ActivateSubscription")
	defer span.End()

	return d.service.ActivateSubscription(ctx, id, correlationID)
}
