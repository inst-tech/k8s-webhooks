package auditor

import (
	"context"
	log "github.com/inst-tech/k8s-webhooks/webhook-reciever/internal/logging"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"k8s.io/apiserver/pkg/apis/audit"
)

var (
	eventListsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auditor_processed_eventlists_total",
		Help: "The total number of processed event lists",
	})
)

var (
	eventsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auditor_processed_events_total",
		Help: "The total number of processed events",
	})
)

func NewAuditor() *processors.EventListProcessor {

	auditor := processors.EventListProcessor{Process: process}
	return &auditor
}

func process(ctx *context.Context, el *audit.EventList) {
	logger := log.GlobalLogger()
	eventListSpan := opentracing.SpanFromContext(*ctx)
	eventListSpan.SetOperationName("handleEventList")
	defer eventListSpan.Finish()
	for _, item := range el.Items {
		eventListContext := opentracing.ContextWithSpan(context.Background(), eventListSpan)
		handleEvent(&eventListContext, &item)
	}
	eventListsProcessed.Inc()
	logger.WithContext(*ctx).Debug("Event List Processed")
}

func handleEvent(ctx *context.Context, e *audit.Event) {
	logger := log.GlobalLogger()
	eventSpan := opentracing.SpanFromContext(*ctx)
	eventSpan.SetOperationName("handleEventList")
	defer eventSpan.Finish()
	eventsProcessed.Inc()
	logger.WithContext(*ctx).WithField("event", &e).Debug("Event Processed")
}
