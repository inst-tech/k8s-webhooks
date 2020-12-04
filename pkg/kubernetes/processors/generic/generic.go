package generic

import (
	"context"
	"encoding/json"
	log "github.com/inst-tech/k8s-webhooks/webhook-reciever/internal/logging"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors/admitter"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors/auditor"
	"github.com/inst-tech/k8s-webhooks/webhook-reciever/pkg/kubernetes/processors/authorizer"
	"github.com/opentracing/opentracing-go"
	"io/ioutil"
	admissionv1 "k8s.io/api/admission/v1beta1"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apiserver/pkg/apis/audit"
	"net/http"
)

const (
	SubjectAccessReview = "SubjectAccessReview"
	AdmissionReview     = "AdmissionReview"
	EventList           = "EventList"
)

func NewGenericProcessor() *processors.GenericProcessor {
	processor := processors.GenericProcessor{Process: ProcessRequest}
	return &processor
}

func ProcessRequest(ctx *context.Context, res http.ResponseWriter, req *http.Request) {
	logger := log.GlobalLogger()
	incomingSpan := opentracing.SpanFromContext(*ctx)
	incomingSpan.SetOperationName("createIncomingRequest")
	defer incomingSpan.Finish()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.WithError(err).Error("ERROR: Unable to read POST body")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	processors.IncomingRequestsProcessed.Inc()

	var incomingRequest processors.IncomingRequest
	err1 := json.Unmarshal(body, &incomingRequest)
	if err1 != nil {
		logger.WithError(err1).WithField("body", string(body[:])).Error("ERROR: Unable to parse POST to JSON")
		return
	}
	logger.WithField("Kind", incomingRequest.Kind).Debug("Processing Incoming Request")

	handlerSpan := opentracing.SpanFromContext(opentracing.ContextWithSpan(context.Background(), incomingSpan))
	handlerSpan.SetOperationName("handleIncomingRequest")
	defer handlerSpan.Finish()

	var response *processors.OutgoingResponse
	var processingError error

	switch incomingRequest.Kind {

	case SubjectAccessReview:
		processors.AuthorizationRequests.Inc()
		subjectAccessReview := authorizationv1.SubjectAccessReview{
			TypeMeta: incomingRequest.TypeMeta,
			Spec:     incomingRequest.Spec,
			Status:   incomingRequest.Status,
		}
		sarProcessor := authorizer.NewAuthorizer()

		response, processingError = sarProcessor.Process(ctx, &subjectAccessReview)
	case AdmissionReview:
		processors.AdmissionRequests.Inc()
		admissionReview := admissionv1.AdmissionReview{
			TypeMeta: incomingRequest.TypeMeta,
			Request:  incomingRequest.Request,
			Response: incomingRequest.Response,
		}
		admissionReviewProcessor := admitter.NewAdmitter()
		response, processingError = admissionReviewProcessor.Process(ctx, &admissionReview)
	case EventList:
		processors.EventListsProcessed.Inc()
		logger.WithField("ItemCount", len(incomingRequest.Items)).Debug("Processing EventList Request")
		eventList := audit.EventList{TypeMeta: incomingRequest.TypeMeta, ListMeta: incomingRequest.ListMeta, Items: incomingRequest.Items}
		eventListProcessor := auditor.NewAuditor()
		eventListProcessor.Process(ctx, &eventList)
		response = &processors.OutgoingResponse{
			TypeMeta:      eventList.TypeMeta,
			AuditResponse: processors.EventListResponse{Status: "OK", StatusMessage: "EventListAccepted"},
		}
	}

	if processingError != nil {
		http.Error(res, processingError.Error(), http.StatusInternalServerError)
		return
	}

	requestResponseSpan := opentracing.SpanFromContext(opentracing.ContextWithSpan(context.Background(), handlerSpan))
	requestResponseSpan.SetOperationName("respondIncomingRequest")
	defer requestResponseSpan.Finish()

	js, err2 := json.Marshal(response)
	if err2 != nil {
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	_, err3 := res.Write(js)
	if err3 != nil {
		logger.WithError(err3).WithField("body", string(body[:])).Error("ERROR: Unable to send response.")
		return
	}

	logger.WithField("request", incomingRequest).Info("Incoming Request Parsed")
}
