package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Config representa as configurações da aplicação
type Config struct {
	Database  DatabaseConfig
	Server    ServerConfig
	Telemetry TelemetryConfig
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// ServerConfig configurações do servidor
type ServerConfig struct {
	Port string
}

// TelemetryConfig configurações de telemetria
type TelemetryConfig struct {
	ServiceName      string
	ServiceVersion   string
	ExporterEndpoint string
}

// LoadConfig carrega as configurações da aplicação
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "payments_customer"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8888"),
		},
		Telemetry: TelemetryConfig{
			ServiceName:      getEnv("SERVICE_NAME", "payments-customer"),
			ServiceVersion:   getEnv("SERVICE_VERSION", "1.0.0"),
			ExporterEndpoint: getEnv("OTEL_EXPORTER_ENDPOINT", "otlcollector:4318"),
		},
	}
}

// NewDatabaseConnection cria uma nova conexão com o banco de dados
func (c *Config) NewDatabaseConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao fazer ping no banco: %w", err)
	}

	return db, nil
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
