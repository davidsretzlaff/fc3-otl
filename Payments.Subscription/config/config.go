package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Config representa as configurações da aplicação
type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Telemetry struct {
		ServiceName    string
		ServiceVersion string
	}
	CustomerServiceURL string
}

// LoadConfig carrega as configurações da aplicação
func LoadConfig() *Config {
	cfg := &Config{}

	// Configurações do servidor
	cfg.Server.Port = getEnvOrDefault("SERVER_PORT", "8081")

	// Configurações do banco de dados
	cfg.Database.Host = getEnvOrDefault("DB_HOST", "localhost")
	cfg.Database.Port = getEnvOrDefault("DB_PORT", "3306")
	cfg.Database.User = getEnvOrDefault("DB_USER", "root")
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", "root")
	cfg.Database.Name = getEnvOrDefault("DB_NAME", "subscription")

	// Configurações de telemetria
	cfg.Telemetry.ServiceName = getEnvOrDefault("TELEMETRY_SERVICE_NAME", "subscription-service")
	cfg.Telemetry.ServiceVersion = getEnvOrDefault("TELEMETRY_SERVICE_VERSION", "1.0.0")

	// URL do serviço de Customer
	cfg.CustomerServiceURL = getEnvOrDefault("CUSTOMER_SERVICE_URL", "http://payments.customer/api/customer")

	return cfg
}

// NewDatabaseConnection cria uma nova conexão com o banco de dados
func (c *Config) NewDatabaseConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
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
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
