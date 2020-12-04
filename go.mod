module github.com/inst-tech/k8s-webhooks/webhook-reciever

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/apiserver v0.19.4
)
