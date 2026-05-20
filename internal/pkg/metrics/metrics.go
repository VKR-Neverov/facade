package metrics

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"strings"
)

var metricsAppName = strings.ReplaceAll(os.Getenv("APP_NAME"), "-", "_")

var reg *prometheus.Registry

var (
	Requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_requests", metricsAppName),
			Help: "Count of requests by handlers",
		},
		[]string{"handler"},
	)
	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_response_time", metricsAppName),
			Help:    "Response time of requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"handler"},
	)
	StatusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_status_codes", metricsAppName),
			Help: "Handlers not successful status codes",
		},
		[]string{"handler", "status_code"},
	)
)

func Init() {
	reg = prometheus.NewRegistry()

	reg.MustRegister(
		Requests,
		ResponseTime,
		StatusCodes,
	)
}

func Run(ctx context.Context, port string) error {
	// Create a new registry
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}
