package processors

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	admissionv1 "k8s.io/api/admission/v1beta1"
	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/apis/audit"
	"net/http"
)

type EventListResponse struct {
	Status string `json:",inline"`
	// +optional
	StatusMessage string `json:"response,omitempty"`
}

type IncomingRequest struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta
	// Request describes the attributes for the admissionv1 request.
	// +optional
	Request *admissionv1.AdmissionRequest `json:"request,omitempty" protobuf:"bytes,1,opt,name=request"`
	// Response describes the attributes for the admissionv1 response.
	// +optional
	Response *admissionv1.AdmissionResponse `json:"response,omitempty" protobuf:"bytes,2,opt,name=response"`
	// Items the items for an Audit EventList
	// +optional
	Items []audit.Event
	// Request describes the attributes for the authorizationv1 request.
	// Spec holds information about the request being evaluated
	Spec authorizationv1.SubjectAccessReviewSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// Response describes the attributes for the authorizationv1 response.
	// Status is filled in by the server and indicates whether the request is allowed or not
	// +optional
	Status authorizationv1.SubjectAccessReviewStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

//goland:noinspection GoVetStructTag
type OutgoingResponse struct {
	metav1.TypeMeta `json:",inline"`
	// Response describes the attributes for the admissionv1 response.
	// +optional
	AdmissionResponse *admissionv1.AdmissionResponse `json:"response,omitempty" protobuf:"bytes,2,opt,name=response"`
	// Response describes the attributes for the audit response.
	// +optional
	AuditResponse EventListResponse `json:"response,omitempty,inline"`
	// Response describes the attributes for the authorizationv1 response.
	//Status is filled in by the server and indicates whether the request is allowed or not
	// +optional
	AccessReviewResponse authorizationv1.SubjectAccessReviewStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ProcessAuthorizationRequest func(ctx *context.Context, sar *authorizationv1.SubjectAccessReview) (res *OutgoingResponse, err error)
type AuthorizationRequestProcessor struct {
	Process ProcessAuthorizationRequest
}

type ProcessAdmissionReview func(ctx *context.Context, ar *admissionv1.AdmissionReview) (res *OutgoingResponse, err error)
type AdmissionReviewProcessor struct {
	Process ProcessAdmissionReview
}

type ProcessEventList func(ctx *context.Context, el *audit.EventList)
type EventListProcessor struct {
	Process ProcessEventList
}

type ProcessIncomingRequest func(ctx *context.Context, res http.ResponseWriter, req *http.Request)
type GenericProcessor struct {
	Process ProcessIncomingRequest
}

//goland:noinspection SpellCheckingInspection
var (
	IncomingRequestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processor_processed_requests_total",
		Help: "The total number of processed incoming requests",
	})
)

//goland:noinspection SpellCheckingInspection
var (
	EventListsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processor_processed_eventlists_total",
		Help: "The total number of processed event lists",
	})
)

//goland:noinspection SpellCheckingInspection
var (
	AuthorizationRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processor_processed_authorizationrequests_total",
		Help: "The total number of processed authorization requests",
	})
)

//goland:noinspection SpellCheckingInspection
var (
	AdmissionRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processor_processed_admissionrequests_total",
		Help: "The total number of processed admission requests",
	})
)
