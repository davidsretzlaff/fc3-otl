package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	}
}

func (c *CustomerClient) CreateCustomer(ctx context.Context, request CustomerRequest) (*CustomerResponse, error) {
	ctx, span := c.tracer.Start(ctx, "CustomerClient.CreateCustomer")
	defer span.End()

	span.SetAttributes(
		attribute.String("customer.name", request.Name),
		attribute.String("customer.email", request.Email),
	)

	url := c.baseURL

	jsonData, err := json.Marshal(request)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao serializar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Propaga o contexto de tracing
	c.propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("erro do servidor: status code %d", resp.StatusCode)
		span.RecordError(err)
		return nil, err
	}

	var customerResponse CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customerResponse); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	span.SetAttributes(
		attribute.String("customer.id", customerResponse.ID),
	)

	return &customerResponse, nil
}
