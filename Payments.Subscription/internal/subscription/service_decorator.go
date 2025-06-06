package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"payments-subscription/internal/common/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	warningThreshold = 2 * time.Second
	errorThreshold   = 5 * time.Second
)

// SubscriptionServiceTracingDecorator é um decorator que adiciona tracing e logging ao serviço
type SubscriptionServiceTracingDecorator struct {
	service SubscriptionServiceInterface
	tracer  trace.Tracer
	logger  *logging.StructuredLogger
}

// NewSubscriptionServiceTracingDecorator cria uma nova instância do decorator
func NewSubscriptionServiceTracingDecorator(service SubscriptionServiceInterface, tracer trace.Tracer) SubscriptionServiceInterface {
	return &SubscriptionServiceTracingDecorator{
		service: service,
		tracer:  tracer,
		logger:  logging.NewStructuredLogger("subscription-service"),
	}
}

// addRequestToSpan adiciona dados da request ao span
func (d *SubscriptionServiceTracingDecorator) addRequestToSpan(span trace.Span, req interface{}) {
	if reqJSON, err := json.Marshal(req); err == nil {
		span.SetAttributes(attribute.String("request", string(reqJSON)))
	}
}

// addResponseToSpan adiciona dados da response ao span
func (d *SubscriptionServiceTracingDecorator) addResponseToSpan(span trace.Span, resp interface{}, err error) {
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return
	}

	if respJSON, err := json.Marshal(resp); err == nil {
		span.SetAttributes(attribute.String("response", string(respJSON)))
	}
}

// logExecutionTime loga o tempo de execução e erros usando logger estruturado
func (d *SubscriptionServiceTracingDecorator) logExecutionTime(ctx context.Context, methodName string, start time.Time, err error) {
	duration := time.Since(start)

	// Log de erro se houver
	if err != nil {
		d.logger.Error(ctx, "MethodExecution",
			fmt.Sprintf("Method %s failed after %v", methodName, duration),
			err,
			map[string]interface{}{
				"method":      methodName,
				"duration_ms": duration.Milliseconds(),
			})
		return
	}

	// Log baseado no tempo de execução
	switch {
	case duration >= errorThreshold:
		d.logger.Error(ctx, "MethodPerformance",
			fmt.Sprintf("Method %s took too long to execute: %v (threshold: %v)", methodName, duration, errorThreshold),
			nil,
			map[string]interface{}{
				"method":       methodName,
				"duration_ms":  duration.Milliseconds(),
				"threshold_ms": errorThreshold.Milliseconds(),
			})
	case duration >= warningThreshold:
		d.logger.Info(ctx, "MethodPerformance",
			fmt.Sprintf("Method %s is running slow: %v (threshold: %v)", methodName, duration, warningThreshold),
			map[string]interface{}{
				"method":       methodName,
				"duration_ms":  duration.Milliseconds(),
				"threshold_ms": warningThreshold.Milliseconds(),
			})
	}
}

// CreateSubscription adiciona tracing e logging à operação de criação
func (d *SubscriptionServiceTracingDecorator) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	start := time.Now()
	ctx, span := d.tracer.Start(ctx, "Service.CreateSubscription")
	defer span.End()

	// Adiciona request ao span
	d.addRequestToSpan(span, req)

	response, err := d.service.CreateSubscription(ctx, req)

	// Adiciona response ou erro ao span
	d.addResponseToSpan(span, response, err)

	d.logExecutionTime(ctx, "CreateSubscription", start, err)
	return response, err
}

// GetSubscriptionByID adiciona tracing e logging à operação de busca por ID
func (d *SubscriptionServiceTracingDecorator) GetSubscriptionByID(ctx context.Context, id string) (*SubscriptionResponse, error) {
	start := time.Now()
	ctx, span := d.tracer.Start(ctx, "Service.GetSubscriptionByID")
	defer span.End()

	// Adiciona request ao span
	span.SetAttributes(attribute.String("subscription_id", id))

	response, err := d.service.GetSubscriptionByID(ctx, id)

	// Adiciona response ou erro ao span
	d.addResponseToSpan(span, response, err)

	d.logExecutionTime(ctx, "GetSubscriptionByID", start, err)
	return response, err
}

// GetAllSubscriptions adiciona tracing e logging à operação de busca de todas as subscriptions
func (d *SubscriptionServiceTracingDecorator) GetAllSubscriptions(ctx context.Context) ([]*SubscriptionResponse, error) {
	start := time.Now()
	ctx, span := d.tracer.Start(ctx, "Service.GetAllSubscriptions")
	defer span.End()

	response, err := d.service.GetAllSubscriptions(ctx)

	// Adiciona response ou erro ao span
	d.addResponseToSpan(span, response, err)

	// Adiciona contagem de resultados se não houver erro
	if err == nil {
		span.SetAttributes(attribute.Int("response.count", len(response)))
	}

	d.logExecutionTime(ctx, "GetAllSubscriptions", start, err)
	return response, err
}

// ActivateSubscription adiciona tracing e logging à operação de ativação
func (d *SubscriptionServiceTracingDecorator) ActivateSubscription(ctx context.Context, id, correlationID string) error {
	start := time.Now()
	ctx, span := d.tracer.Start(ctx, "Service.ActivateSubscription")
	defer span.End()

	// Adiciona dados da request ao span
	span.SetAttributes(
		attribute.String("subscription_id", id),
		attribute.String("correlation_id", correlationID),
	)

	err := d.service.ActivateSubscription(ctx, id, correlationID)

	// Adiciona erro ao span se houver
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
	}

	d.logExecutionTime(ctx, "ActivateSubscription", start, err)
	return err
}
