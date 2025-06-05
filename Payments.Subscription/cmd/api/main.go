package main

import (
	"log"
	"net/http"

	"payments-subscription/config"
	"payments-subscription/internal/common/logging"
	"payments-subscription/internal/common/middleware"
	opentel "payments-subscription/internal/common/telemetry"
	"payments-subscription/internal/customer"
	"payments-subscription/internal/subscription"
	mysql "payments-subscription/internal/subscription/mysql"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {
	// Carrega as configurações
	cfg := config.LoadConfig()

	// Inicializa o logger estruturado
	logger := logging.NewStructuredLogger("subscription-service")

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

	// Cria o cliente do serviço de Customer
	customerClient := customer.NewCustomerClient(cfg.CustomerServiceURL)

	// Cria o serviço base
	subscriptionService := subscription.NewSubscriptionService(repositoryDecored, subscriptionEventService, customerClient)

	// Aplica o decorador de tracing
	subscriptionServiceDecored := subscription.NewSubscriptionServiceTracingDecorator(subscriptionService, tracer)

	subscriptionHandler := subscription.NewSubscriptionHandler(subscriptionServiceDecored)

	// Configura o router HTTP com middleware de tracing
	router := mux.NewRouter()

	// Middlewares em ordem:
	// 1. Correlation ID (primeiro para garantir que todas as operações tenham correlation ID)
	router.Use(middleware.CorrelationIDMiddleware)

	// 2. Logging HTTP (depois do correlation ID para logar com o ID correto)
	router.Use(middleware.LoggingMiddleware(logger))

	// 3. OpenTelemetry tracing
	router.Use(otelmux.Middleware(ot.ServiceName,
		otelmux.WithTracerProvider(ot.GetTracerProvider()),
		otelmux.WithPropagators(ot.GetPropagators()),
	))

	// Configura as rotas
	router.HandleFunc("/subscriptions", subscriptionHandler.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions/{id}", subscriptionHandler.GetSubscriptionByID).Methods("GET")
	router.HandleFunc("/subscriptions", subscriptionHandler.GetAllSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions/{id}/activate", subscriptionHandler.ActivateSubscription).Methods("POST")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Info(ctx, "HEALTH_CHECK", "Health check acessado", nil)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Log de inicialização
	log.Printf("Subscription Service iniciado na porta %s", cfg.Server.Port)
	log.Printf("Customer Service URL: %s", cfg.CustomerServiceURL)

	// Inicia o servidor
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
