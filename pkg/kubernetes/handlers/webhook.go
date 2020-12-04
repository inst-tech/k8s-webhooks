package handlers

import (
	"context"
	log "github.com/inst-tech/k8s-webhooks/webhook-reciever/internal/logging"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

var logger log.StandardLogger

type RequestHandler struct {
	processor *processors.GenericProcessor
}

var requestHandler RequestHandler

func NewRequestHandler(processor *processors.GenericProcessor) RequestHandler {
	logger = log.GlobalLogger()
	requestHandler = RequestHandler{processor: processor}
	return requestHandler
}

func (r *RequestHandler) Handle(res http.ResponseWriter, req *http.Request) {
	var serverSpan opentracing.Span

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	if err != nil {
		logger.WithError(err).Error("ERROR: Unable to extract Opentracing Headers from Request")
	}

	serverSpan = opentracing.StartSpan(
		"processIncomingRequest",
		ext.RPCServerOption(wireContext))
	defer serverSpan.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)
	requestHandler.processor.Process(&ctx, res, req)
}
