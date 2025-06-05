package logging

import (
	"context"
	"fmt"
	"log"
	"time"
)

// LogLevel representa o nível do log
type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
)

// StructuredLogger é um logger estruturado simples
type StructuredLogger struct {
	serviceName string
}

// NewStructuredLogger cria um novo logger estruturado
func NewStructuredLogger(serviceName string) *StructuredLogger {
	return &StructuredLogger{
		serviceName: serviceName,
	}
}

// formatMessage formata mensagem de log de forma limpa
func (l *StructuredLogger) formatMessage(ctx context.Context, level LogLevel, message string) string {
	timestamp := time.Now().UTC().Format("15:04:05")

	// Extrair correlation ID do contexto se existir
	if corrID := GetCorrelationID(ctx); corrID != "" {
		if level == ERROR {
			return fmt.Sprintf("%s [subscription] [CorrelationId:%s] ERROR: %s", timestamp, corrID, message)
		}
		return fmt.Sprintf("%s [subscription] [CorrelationId:%s] %s", timestamp, corrID, message)
	}

	// Log sem correlation ID
	if level == ERROR {
		return fmt.Sprintf("%s [subscription] ERROR: %s", timestamp, message)
	}
	return fmt.Sprintf("%s [subscription] %s", timestamp, message)
}

// Info registra um log de nível INFO - APENAS ESSENCIAIS
func (l *StructuredLogger) Info(ctx context.Context, operation, message string, contextData map[string]interface{}) {
	formattedMessage := l.formatMessage(ctx, INFO, message)
	log.Println(formattedMessage)
}

// Error registra um log de nível ERROR - APENAS ESSENCIAIS
func (l *StructuredLogger) Error(ctx context.Context, operation, message string, err error, contextData map[string]interface{}) {
	errorMessage := message
	if err != nil {
		errorMessage = fmt.Sprintf("%s: %s", message, err.Error())
	}

	formattedMessage := l.formatMessage(ctx, ERROR, errorMessage)
	log.Println(formattedMessage)
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
