package main

import (
	"log"
	"net/http"

	"payments-customer/config"
	opentel "payments-customer/internal/common/telemetry"
	"payments-customer/internal/subscription"
	mysql "payments-customer/internal/subscription/mysql"

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

	// Obtém o tracer configurado
	tracer := ot.GetTracer()

	// Conecta com o banco de dados
	db, err := cfg.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco: %v", err)
	}
	defer db.Close()

	// Inicializa as dependências seguindo DDD
	repository := mysql.NewMySQLSubscriptionRepository(db)

	// Aplica o decorator de tracing ao repositório
	repositoryDecored := subscription.NewSubscriptionRepositoryTracingDecorator(repository, tracer)

	subscriptionEventPublisher := subscription.NewInMemoryEventPublisher()
	subscriptionEventService := subscription.NewSubscriptionEventService(subscriptionEventPublisher)

	// Cria o serviço base
	subscriptionService := subscription.NewSubscriptionService(repositoryDecored, subscriptionEventService)

	// Aplica o decorador de tracing
	subscriptionServiceDecored := subscription.NewSubscriptionServiceTracingDecorator(subscriptionService, tracer)

	subscriptionHandler := subscription.NewSubscriptionHandler(subscriptionServiceDecored)

	// Configura o router HTTP com middleware de tracing
	router := mux.NewRouter()
	router.Use(otelmux.Middleware(ot.ServiceName,
		otelmux.WithTracerProvider(ot.GetTracerProvider()),
		otelmux.WithPropagators(ot.GetPropagators()),
	))

	// Registra as rotas
	subscriptionHandler.RegisterRoutes(router)

	// Adiciona uma rota de health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Printf("Servidor iniciado na porta %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}
