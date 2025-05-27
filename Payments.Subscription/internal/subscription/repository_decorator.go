package subscription

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// SubscriptionRepositoryTracingDecorator é um decorator que adiciona tracing ao repositório
type SubscriptionRepositoryTracingDecorator struct {
	repository SubscriptionRepository
	tracer     trace.Tracer
}

// NewSubscriptionRepositoryTracingDecorator cria uma nova instância do decorator de tracing
func NewSubscriptionRepositoryTracingDecorator(repository SubscriptionRepository, tracer trace.Tracer) SubscriptionRepository {
	return &SubscriptionRepositoryTracingDecorator{
		repository: repository,
		tracer:     tracer,
	}
}

// Create adiciona tracing à operação de criação
func (d *SubscriptionRepositoryTracingDecorator) Create(ctx context.Context, subscription *Subscription) error {
	ctx, span := d.tracer.Start(ctx, "Repository.Create")
	defer span.End()

	return d.repository.Create(ctx, subscription)
}

// GetByID adiciona tracing à operação de busca por ID
func (d *SubscriptionRepositoryTracingDecorator) GetByID(ctx context.Context, id SubscriptionID) (*Subscription, error) {
	ctx, span := d.tracer.Start(ctx, "Repository.GetByID")
	defer span.End()

	return d.repository.GetByID(ctx, id)
}

// GetByCustomerID adiciona tracing à operação de busca por customer ID
func (d *SubscriptionRepositoryTracingDecorator) GetByCustomerID(ctx context.Context, customerID CustomerID) ([]*Subscription, error) {
	ctx, span := d.tracer.Start(ctx, "Repository.GetByCustomerID")
	defer span.End()

	return d.repository.GetByCustomerID(ctx, customerID)
}

// Update adiciona tracing à operação de atualização
func (d *SubscriptionRepositoryTracingDecorator) Update(ctx context.Context, subscription *Subscription) error {
	ctx, span := d.tracer.Start(ctx, "Repository.Update")
	defer span.End()

	return d.repository.Update(ctx, subscription)
}

// GetAll adiciona tracing à operação de busca de todas as subscriptions
func (d *SubscriptionRepositoryTracingDecorator) GetAll(ctx context.Context) ([]*Subscription, error) {
	ctx, span := d.tracer.Start(ctx, "Repository.GetAll")
	defer span.End()

	return d.repository.GetAll(ctx)
}
