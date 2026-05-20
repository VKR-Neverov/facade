package config

import (
	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"os"
	"time"
)

var Tracer opentracing.Tracer

// New for current service
func New() (tracer opentracing.Tracer, closer io.Closer, err error) {
	jaegerEndpoint := "http://jaeger:14268/api/traces"

	serviceName := os.Getenv("APP_NAME")

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   jaegerEndpoint,
		},
		ServiceName: serviceName,
	}

	tracer, closer, err = cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	Tracer = tracer

	opentracing.SetGlobalTracer(Tracer)
	return
}
