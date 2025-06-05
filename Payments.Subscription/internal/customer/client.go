package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"payments-subscription/internal/common/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type CustomerClient struct {
	baseURL    string
	httpClient *http.Client
	propagator propagation.TextMapPropagator
	tracer     trace.Tracer
	logger     *logging.StructuredLogger
}

type CustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewCustomerClient(baseURL string) *CustomerClient {
	return &CustomerClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		propagator: otel.GetTextMapPropagator(),
		tracer:     otel.GetTracerProvider().Tracer("customer-client"),
		logger:     logging.NewStructuredLogger("subscription-service"),
	}
}

func (c *CustomerClient) CreateCustomer(ctx context.Context, request CustomerRequest) (*CustomerResponse, error) {
	startTime := time.Now()
	operation := "CustomerClient.CreateCustomer"

	// Garantir que existe correlation ID
	ctx = logging.EnsureCorrelationID(ctx, "subscription")
	correlationID := logging.GetCorrelationID(ctx)

	// Log início da operação
	c.logger.OperationStart(ctx, operation, map[string]interface{}{
		"customer_email": request.Email,
		"customer_name":  request.Name,
	})

	ctx, span := c.tracer.Start(ctx, operation)
	defer span.End()

	span.SetAttributes(
		attribute.String("customer.name", request.Name),
		attribute.String("customer.email", request.Email),
		attribute.String("correlation.id", correlationID),
	)

	url := c.baseURL

	jsonData, err := json.Marshal(request)
	if err != nil {
		span.RecordError(err)
		c.logger.Error(ctx, operation, "Failed to serialize request", err, nil)
		return nil, fmt.Errorf("erro ao serializar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		span.RecordError(err)
		c.logger.Error(ctx, operation, "Failed to create HTTP request", err, nil)
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", correlationID)

	// Propagar contexto de tracing via headers W3C
	c.propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := c.httpClient.Do(req)

	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
		defer resp.Body.Close()
	}

	if err != nil {
		c.logger.LogServiceCall(ctx, "Customer", statusCode, err)
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao fazer request: %w", err)
	}

	// Log para status codes de erro
	if statusCode != http.StatusCreated && statusCode != http.StatusOK {
		err := fmt.Errorf("customer service returned status code %d", statusCode)
		c.logger.LogServiceCall(ctx, "Customer", statusCode, nil)
		span.RecordError(err)
		return nil, err
	}

	var customerResponse CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customerResponse); err != nil {
		span.RecordError(err)
		c.logger.Error(ctx, operation, "Failed to decode response", err, nil)
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	span.SetAttributes(
		attribute.String("customer.id", customerResponse.ID),
	)

	// Log fim apenas se demorou muito
	c.logger.OperationEnd(ctx, operation, startTime, map[string]interface{}{
		"customer_id": customerResponse.ID,
	})

	return &customerResponse, nil
}
