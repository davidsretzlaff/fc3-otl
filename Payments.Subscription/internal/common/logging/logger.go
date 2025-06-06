package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LogLevel representa o nível do log
type LogLevel string

const (
	INFO  LogLevel = "info"
	ERROR LogLevel = "error"
	WARN  LogLevel = "warn"
)

// LogEntry representa a estrutura do log JSON
type LogEntry struct {
	Time          string `json:"time"`
	Level         string `json:"level"`
	Message       string `json:"msg"`
	CorrelationID string `json:"correlation_id"`
	Service       string `json:"service"`
}

// StructuredLogger é um logger estruturado que gera JSON
type StructuredLogger struct {
	serviceName string
}

// NewStructuredLogger cria um novo logger estruturado
func NewStructuredLogger(serviceName string) *StructuredLogger {
	return &StructuredLogger{
		serviceName: serviceName,
	}
}

// logJSON registra uma entrada de log em formato JSON
func (l *StructuredLogger) logJSON(ctx context.Context, level LogLevel, message string) {
	entry := LogEntry{
		Time:          time.Now().UTC().Format(time.RFC3339),
		Level:         string(level),
		Message:       message,
		CorrelationID: GetCorrelationID(ctx),
		Service:       l.serviceName,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback para stderr se houver erro no JSON
		fmt.Fprintf(os.Stderr, "ERROR: Failed to marshal log entry: %v\n", err)
		return
	}

	// Escrever JSON diretamente no stdout sem timestamp adicional
	fmt.Println(string(jsonBytes))
}

// Info registra um log de nível INFO - APENAS ESSENCIAIS
func (l *StructuredLogger) Info(ctx context.Context, operation, message string, contextData map[string]interface{}) {
	l.logJSON(ctx, INFO, message)
}

// Error registra um log de nível ERROR - APENAS ESSENCIAIS
func (l *StructuredLogger) Error(ctx context.Context, operation, message string, err error, contextData map[string]interface{}) {
	errorMessage := message
	if err != nil {
		errorMessage = fmt.Sprintf("%s: %s", message, err.Error())
	}

	l.logJSON(ctx, ERROR, errorMessage)
}

// OperationStart - LOG ESSENCIAL apenas para operações principais
func (l *StructuredLogger) OperationStart(ctx context.Context, operation string, contextData map[string]interface{}) {
	// Apenas para operações principais de negócio
	if operation == "CreateSubscription" {
		if email, ok := contextData["customer_email"].(string); ok {
			message := fmt.Sprintf("Starting CreateSubscription for %s", email)
			l.Info(ctx, operation, message, nil)
		}
	}

	if operation == "CustomerClient.CreateCustomer" {
		message := "Calling customer service to create customer"
		l.Info(ctx, operation, message, nil)
	}
}

// OperationEnd - LOG ESSENCIAL apenas para erros ou operações muito lentas
func (l *StructuredLogger) OperationEnd(ctx context.Context, operation string, startTime time.Time, contextData map[string]interface{}) {
	duration := time.Since(startTime)

	// Log apenas se for muito lento (>3 segundos)
	if duration.Milliseconds() > 3000 {
		message := fmt.Sprintf("SLOW operation %s completed in %dms", operation, duration.Milliseconds())
		l.Info(ctx, operation, message, nil)
	}
}

// LogServiceCall - LOG para chamadas entre serviços
func (l *StructuredLogger) LogServiceCall(ctx context.Context, service string, statusCode int, err error) {
	if err != nil {
		message := fmt.Sprintf("%s service call failed", service)
		l.Error(ctx, "ServiceCall", message, err, nil)
		return
	}

	if statusCode >= 400 {
		message := fmt.Sprintf("%s service returned status code %d", service, statusCode)
		l.Info(ctx, "ServiceCall", message, nil)
	}
}
