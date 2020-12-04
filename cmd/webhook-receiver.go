package main

import (
	"github.com/gorilla/mux"
	log "github.com/inst-tech/k8s-webhooks/webhook-reciever/internal/logging"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/handlers"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors/generic"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"net/http"
	"time"
)

const ApplicationName = "webhookreceiver"

var logger *log.StandardLogger

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: ApplicationName + "_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

func main() {
	logger = log.NewLogger(ApplicationName)
	log.SetGlobalLogger(*logger)

	cfg, err := jaegercfg.FromEnv()
	if err != nil {

		return
	}
	cfg.ServiceName = ApplicationName

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(logger),
	)

	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)

	if err != nil {
		logger.WithError(err).Error("Could not initialize Tracing")
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer closer.Close()

	processor := generic.NewGenericProcessor()
	requestHandler := handlers.NewRequestHandler(processor)

	r := mux.NewRouter()

	r.Use(prometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/audit", requestHandler.Handle)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err1 := srv.ListenAndServe()
	if err1 != nil {
		logger.WithError(err).Error("Could not ListenAndServe")
	}
}
