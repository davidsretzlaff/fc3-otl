package logging

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

// Chaves para o contexto
type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
	userIDKey        contextKey = "user_id"
	requestIDKey     contextKey = "request_id"
)

// GenerateCorrelationID gera um correlation ID único
func GenerateCorrelationID(service string) string {
	// Formato: {service}-{timestamp}-{random}
	timestamp := time.Now().Format("20060102150405")

	// Gerar 4 bytes aleatórios
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := fmt.Sprintf("%x", randomBytes)

	return fmt.Sprintf("%s-%s-%s", service, timestamp, randomHex)
}

// WithCorrelationID adiciona correlation ID ao contexto
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

// GetCorrelationID obtém o correlation ID do contexto
func GetCorrelationID(ctx context.Context) string {
	if value := ctx.Value(correlationIDKey); value != nil {
		if corrID, ok := value.(string); ok {
			return corrID
		}
	}
	return ""
}

// WithUserID adiciona user ID ao contexto
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID obtém o user ID do contexto
func GetUserID(ctx context.Context) string {
	if value := ctx.Value(userIDKey); value != nil {
		if userID, ok := value.(string); ok {
			return userID
		}
	}
	return ""
}

// WithRequestID adiciona request ID ao contexto
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestID obtém o request ID do contexto
func GetRequestID(ctx context.Context) string {
	if value := ctx.Value(requestIDKey); value != nil {
		if reqID, ok := value.(string); ok {
			return reqID
		}
	}
	return ""
}

// EnsureCorrelationID garante que existe um correlation ID no contexto
func EnsureCorrelationID(ctx context.Context, service string) context.Context {
	if GetCorrelationID(ctx) == "" {
		corrID := GenerateCorrelationID(service)
		return WithCorrelationID(ctx, corrID)
	}
	return ctx
}
