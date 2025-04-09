package main

import (
	"io"       // Para ler o corpo das respostas HTTP
	"log"      // Para logging
	"net/http" // Para fazer requisições HTTP
	"time"     // Para simular delays

	// Pacotes do OpenTelemetry
	"github.com/codeedu/otel-go/infra/opentel"                                   // Sua implementação do OpenTelemetry
	"github.com/gorilla/mux"                                                     // Router HTTP
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux" // Middleware de tracing para o Mux
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"              // Instrumentação para chamadas HTTP
	"go.opentelemetry.io/otel/baggage"                                           // Para propagar contexto entre serviços
	"go.opentelemetry.io/otel/trace"                                             // Para trabalhar com traces
)

// Variável global para o tracer
var tracer trace.Tracer

func main() {
	// Inicializa o OpenTelemetry
	ot := opentel.NewOpenTel()
	ot.ServiceName = "GoApp"            // Nome do serviço que aparecerá no Jaeger
	ot.ServiceVersion = "0.1"           // Versão do serviço
	ot.ExporterEndpoint = "jaeger:4318" // Endpoint do Jaeger para enviar os traces

	// Obtém o tracer configurado
	tracer = ot.GetTracer()

	// Configura o router HTTP com middleware de tracing
	router := mux.NewRouter()
	router.Use(otelmux.Middleware(ot.ServiceName,
		otelmux.WithTracerProvider(ot.GetTracerProvider()), // Usa o tracer provider configurado
		otelmux.WithPropagators(ot.GetPropagators()),       // Usa os propagadores configurados
	))
	router.HandleFunc("/", homeHandler)  // Rota principal
	http.ListenAndServe(":8888", router) // Inicia o servidor na porta 8888
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	// Cria um contexto sem baggage (informações adicionais)
	ctx := baggage.ContextWithoutBaggage(request.Context())

	// Primeiro span: Processamento inicial
	ctx, initialProcessing := tracer.Start(ctx, "app.initial.processing")
	time.Sleep(time.Millisecond * 100) // Simula processamento
	initialProcessing.End()            // Finaliza o span

	// Segundo span: Chamada para o serviço de dados
	ctx, dataServiceCall := tracer.Start(ctx, "app.data.service.call")
	// Configura o cliente HTTP com instrumentação de tracing
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	// Cria a requisição com o contexto de tracing
	req, err := http.NewRequestWithContext(ctx, "GET", "http://netcoreapp:80/MyController", nil)
	if err != nil {
		log.Fatal(err)
	}
	// Faz a requisição
	res, err := client.Do(req)

	if err != nil {
		log.Printf("Erro ao fazer requisição: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close() // Garante que o body será fechado

	// Lê o corpo da resposta
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Erro ao ler resposta: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	time.Sleep(time.Millisecond * 300) // Simula processamento
	dataServiceCall.End()              // Finaliza o span da chamada HTTP

	// Terceiro span: Preparação da resposta
	_, responsePreparation := tracer.Start(ctx, "app.response.preparation")
	time.Sleep(time.Millisecond * 100) // Simula processamento
	writer.WriteHeader(http.StatusOK)  // Define status 200
	writer.Write(body)                 // Escreve o corpo da resposta
	responsePreparation.End()          // Finaliza o span
}
