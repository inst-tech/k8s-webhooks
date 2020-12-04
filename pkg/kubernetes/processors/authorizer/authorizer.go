package authorizer

import (
	"context"
	log "github.com/inst-tech/k8s-webhooks/webhook-reciever/internal/logging"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	authorizationv1 "k8s.io/api/authorization/v1"
)

func NewAuthorizer() *processors.AuthorizationRequestProcessor {
	a := processors.AuthorizationRequestProcessor{Process: process}
	return &a
}

var (
	authorizationRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authorizer_processed_requests_total",
		Help: "The total number of processed authorization requests",
	})
)

var (
	authorizationRequestsAllowed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authorizer_processed_requests_allowed",
		Help: "The total number of allowed authorization requests",
	})
)

var (
	authorizationRequestsDenied = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authorizer_processed_requests_denied",
		Help: "The total number of denied authorization requests",
	})
)

func process(ctx *context.Context, sar *authorizationv1.SubjectAccessReview) (res *processors.OutgoingResponse, err error) {
	logger := log.GlobalLogger()
	sarSpan := opentracing.SpanFromContext(*ctx)
	sarSpan.SetOperationName("handleSubjectAccessReview")
	defer sarSpan.Finish()
	authorizationRequests.Inc()

	logger.WithContext(*ctx).WithField("subjectAccessReview", sar).Debug("Authorization Request Processed")

	sar.Status.Allowed = true
	sar.Status.Reason = "defaultAllow"
	authorizationRequestsAllowed.Inc()

	return &processors.OutgoingResponse{
		TypeMeta:             sar.TypeMeta,
		AccessReviewResponse: sar.Status,
	}, nil

}
