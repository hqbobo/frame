package service

import (
	"context"
	"fmt"

	"github.com/hqbobo/frame/common/log"
	"github.com/micro/go-micro/v2/server"

	goTracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func newTracer(name, url string) goTracing.Tracer {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter(url)
	//defer reporter.Close()
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(name, "")
	if err != nil {
		log.Errorf("unable to create local endpoint: %+v\n", err)
	}
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Errorf("unable to create tracer: %+v\n", err)
	}
	// use zipkin-go-opentracing to wrap our tracer
	return zipkinot.Wrap(nativeTracer)
}

// Span Span结构
type Span interface {
	Finish()
}

type nullTracerSpan struct {
}

func (ns nullTracerSpan) Finish() {

}

// TraceSpan 创建一个链路追踪节点
func TraceSpan(ctx context.Context, name string) (Span, context.Context) {
	if goTracing.IsGlobalTracerRegistered() {
		return goTracing.StartSpanFromContext(ctx, name)
	}
	return nullTracerSpan{}, ctx
}

// PanicHandlerWrapper accepts an opentracing Tracer and returns a Handler Wrapper
func PanicHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			log.Info(name)
			return h(ctx, req, rsp)
		}
	}
}
