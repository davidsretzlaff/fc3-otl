package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomerClient struct {
	baseURL    string
	httpClient *http.Client
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
	}
}

func (c *CustomerClient) CreateCustomer(ctx context.Context, request CustomerRequest) (*CustomerResponse, error) {
	url := c.baseURL

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro do servidor: status code %d", resp.StatusCode)
	}

	var customerResponse CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customerResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &customerResponse, nil
}
