package admitter

import (
	"context"
	"encoding/json"
	log "github.com/inst-tech/k8s-webhooks/webhook-reciever/internal/logging"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	admissionv1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

var (
	admissionRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "admitter_processed_requests_total",
		Help: "The total number of processed admission requests",
	})
)

var (
	admissionRequestsAllowed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "admitter_processed_requests_allowed",
		Help: "The total number of processed admission requests",
	})
)

var (
	admissionRequestsDenied = promauto.NewCounter(prometheus.CounterOpts{
		Name: "admitter_processed_requests_denied",
		Help: "The total number of processed admission requests",
	})
)

func NewAdmitter() *processors.AdmissionReviewProcessor {
	admitter := processors.AdmissionReviewProcessor{Process: process}
	return &admitter
}

func process(ctx *context.Context, ar *admissionv1.AdmissionReview) (res *processors.OutgoingResponse, err error) {
	logger := log.GlobalLogger()
	arSpan := opentracing.SpanFromContext(*ctx)
	arSpan.SetOperationName("handleAdmissionReview")
	defer arSpan.Finish()

	admissionRequests.Inc()

	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		logger.WithError(err).Error("Could not unmarshal raw object: %v")
		ar.Response.Allowed = false
		admissionRequestsDenied.Inc()
		return &processors.OutgoingResponse{
			TypeMeta:          ar.TypeMeta,
			AdmissionResponse: ar.Response,
		}, err
	}

	logger.WithField("kind", req.Kind).WithField("Namespace", req.Namespace).WithField("Name", pod.Name).
		WithField("UID", req.UID).WithField("UserInfo", req.UserInfo).Info("AdmissionReview Complete")

	logger.WithContext(*ctx).WithField("AdmissionRequest", ar.Request).WithField("AdmissionResponse", ar.Response).Debug("Admission Review Processed")
	admissionRequestsAllowed.Inc()

	return &processors.OutgoingResponse{
		TypeMeta:          ar.TypeMeta,
		AdmissionResponse: ar.Response,
	}, nil
}
