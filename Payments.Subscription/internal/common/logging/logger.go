package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	logFile     *os.File
	writers     []io.Writer
}

// NewStructuredLogger cria um novo logger estruturado
func NewStructuredLogger(serviceName string) *StructuredLogger {
	logger := &StructuredLogger{
		serviceName: serviceName,
	}

	// Criar diretório de logs se não existir
	logDir := "/app/logs/apps"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create log directory: %v\n", err)
		// Fallback para stdout se não conseguir criar diretório
		logger.writers = []io.Writer{os.Stdout}
		return logger
	}

	// Criar arquivo de log com data atual
	logFileName := fmt.Sprintf("%s%s.log", serviceName, time.Now().Format("20060102"))
	logFilePath := filepath.Join(logDir, logFileName)

	// Debug: imprimir o caminho do arquivo
	fmt.Fprintf(os.Stderr, "DEBUG: Tentando criar arquivo de log: %s\n", logFilePath)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to open log file %s: %v\n", logFilePath, err)
		// Se falhar ao abrir arquivo, usar apenas stdout
		logger.writers = []io.Writer{os.Stdout}
	} else {
		logger.logFile = logFile
		fmt.Fprintf(os.Stderr, "SUCCESS: Log file created: %s\n", logFilePath)
		// Escrever APENAS no arquivo (não stdout)
		logger.writers = []io.Writer{logFile}
	}

	return logger
}

// Close fecha o arquivo de log
func (l *StructuredLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// logJSON registra uma entrada de log em formato JSON
func (l *StructuredLogger) logJSON(ctx context.Context, level LogLevel, message string) {
	entry := LogEntry{
		Time:          time.Now().UTC().Format(time.RFC3339Nano),
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

	// Escrever em todos os writers (stdout + arquivo)
	jsonString := string(jsonBytes)
	for _, writer := range l.writers {
		fmt.Fprintln(writer, jsonString)
	}
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
			message := fmt.Sprintf("[subscription] Starting CreateSubscription for %s", email)
			l.Info(ctx, operation, message, nil)
		}
	}

	if operation == "CustomerClient.CreateCustomer" {
		message := "[subscription] Calling customer service to create customer"
		l.Info(ctx, operation, message, nil)
	}
}

// OperationEnd - LOG ESSENCIAL apenas para erros ou operações muito lentas
func (l *StructuredLogger) OperationEnd(ctx context.Context, operation string, startTime time.Time, contextData map[string]interface{}) {
	duration := time.Since(startTime)

	// Log apenas se for muito lento (>3 segundos)
	if duration.Milliseconds() > 3000 {
		message := fmt.Sprintf("[subscription] SLOW operation %s completed in %dms", operation, duration.Milliseconds())
		l.Info(ctx, operation, message, nil)
	}
}

// LogServiceCall - LOG para chamadas entre serviços
func (l *StructuredLogger) LogServiceCall(ctx context.Context, service string, statusCode int, err error) {
	if err != nil {
		message := fmt.Sprintf("[subscription] %s service call failed", service)
		l.Error(ctx, "ServiceCall", message, err, nil)
		return
	}

	if statusCode >= 400 {
		message := fmt.Sprintf("[subscription] %s service returned status code %d", service, statusCode)
		l.Info(ctx, "ServiceCall", message, nil)
	}
}
