package main

import (
	"log"
	"net/http"

	"payments-customer/config"
	opentel "payments-customer/internal/common/telemetry"
	"payments-customer/internal/customer"
	mysql "payments-customer/internal/customer/mysql"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {
	// Carrega as configurações
	cfg := config.LoadConfig()

	// Inicializa o OpenTelemetry
	ot := opentel.NewOpenTel()
	ot.ServiceName = cfg.Telemetry.ServiceName
	ot.ServiceVersion = cfg.Telemetry.ServiceVersion
	ot.ExporterEndpoint = cfg.Telemetry.ExporterEndpoint

	// Obtém o tracer configurado
	tracer := ot.GetTracer()

	// Conecta com o banco de dados
	db, err := cfg.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco: %v", err)
	}
	defer db.Close()

	// Inicializa as dependências seguindo DDD
	customerRepository := mysql.NewMySQLCustomerRepository(db, tracer)
	customerService := customer.NewCustomerService(customerRepository, tracer)
	customerHandler := customer.NewCustomerHandler(customerService, tracer)

	// Configura o router HTTP com middleware de tracing
	router := mux.NewRouter()
	router.Use(otelmux.Middleware(ot.ServiceName,
		otelmux.WithTracerProvider(ot.GetTracerProvider()),
		otelmux.WithPropagators(ot.GetPropagators()),
	))

	// Registra as rotas
	customerHandler.RegisterRoutes(router)

	// Adiciona uma rota de health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Printf("Servidor iniciado na porta %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}
