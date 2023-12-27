package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/JailtonJunior94/telemetry-poc/infra/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	otelTelemetry := telemetry.NewOpenTel("GoApp", "1.0.0", "http://localhost:9411/api/v2/spans")
	tracer = otelTelemetry.GetTracer()

	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))
	router.Use(otelchi.Middleware("goapp-server", otelchi.WithChiRoutes(router)))

	router.Get("/", handle)

	server := http.Server{
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "6000"))
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := baggage.ContextWithoutBaggage(r.Context())

	ctx, httpCall := tracer.Start(ctx, "request-get-address")
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/address/06503015", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	httpCall.End()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
