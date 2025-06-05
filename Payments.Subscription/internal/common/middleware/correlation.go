package middleware

import (
	"net/http"

	"payments-subscription/internal/common/logging"
)

// CorrelationIDMiddleware middleware para gerenciar correlation ID
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Tenta extrair correlation ID do header
		correlationID := r.Header.Get("X-Correlation-ID")

		// Se não existir, gera um novo
		if correlationID == "" {
			correlationID = logging.GenerateCorrelationID("subscription")
		}

		// Adiciona o correlation ID ao contexto
		ctx = logging.WithCorrelationID(ctx, correlationID)

		// Adiciona o correlation ID ao header de resposta para facilitar debugging
		w.Header().Set("X-Correlation-ID", correlationID)

		// Continua com o request usando o contexto atualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware middleware SILENCIOSO - remove logs redundantes de HTTP
func LoggingMiddleware(logger *logging.StructuredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// EXECUÇÃO SILENCIOSA - sem logs de requisição HTTP
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wrapper para capturar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
