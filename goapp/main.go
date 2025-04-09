package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/codeedu/otel-go/infra/opentel"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
)

var tracer trace.Tracer

func main() {
	ot := opentel.NewOpenTel()
	ot.ServiceName = "GoApp"
	ot.ServiceVersion = "0.1"
	ot.ExporterEndpoint = "jaeger:4318"

	tracer = ot.GetTracer()

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(ot.ServiceName, 
		otelmux.WithTracerProvider(ot.GetTracerProvider()),
		otelmux.WithPropagators(ot.GetPropagators()),
	))
	router.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8888", router)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := baggage.ContextWithoutBaggage(request.Context())

	// rotina 1 - Process File
	ctx, processFile := tracer.Start(ctx, "process-file")
	time.Sleep(time.Millisecond * 100)
	processFile.End()

	// rotina 2 - Fazer Request HTTP
	ctx, httpCall := tracer.Start(ctx, "request-remote-json")
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://netcoreapp:80/MyController", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req) // chamo a requisição

	if err != nil {
		log.Printf("Erro ao fazer requisição: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Erro ao ler resposta: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	time.Sleep(time.Millisecond * 300)
	httpCall.End()

	// rotina 3 - Exibir resultado
	ctx, renderContent := tracer.Start(ctx, "render-content")
	time.Sleep(time.Millisecond * 100)
	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
	renderContent.End()
}
