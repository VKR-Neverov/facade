package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fidesy-pay/facade/internal/pkg/metrics"
	auth_service "github.com/fidesy-pay/facade/pkg/auth-service"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc/metadata"
)

var Tracer opentracing.Tracer

func Auth(authClient auth_service.AuthServiceClient) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// no need auth for rest api
			if strings.Contains(r.URL.String(), "/api/invoice") {
				methodName := ""
				if r.Method == http.MethodGet {
					methodName = "GET /api/invoice"
				} else {
					methodName = "POST /api/invoice"
				}

				metrics.Requests.WithLabelValues(methodName).Inc()
				defer metrics.ResponseTime.WithLabelValues(methodName).Observe(float64(time.Since(start).Milliseconds()))

				next.ServeHTTP(w, r)
				return
			}

			body := GraphQLRequestBody{}
			data, _ := io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(data))

			if string(data) == "" {
				next.ServeHTTP(w, r)
				return
			}

			err := json.Unmarshal(data, &body)
			if err != nil {
				log.Printf("json.Unmarhsal: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if strings.Contains(*body.Query, "IntrospectionQuery") {
				next.ServeHTTP(w, r)
				return
			}

			if body.OperationName == nil {
				log.Println("black swan, operation name is nil")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			span := Tracer.StartSpan(*body.OperationName)
			defer span.Finish()

			metrics.Requests.WithLabelValues(*body.OperationName).Inc()

			jaegerSpan, _ := span.(*jaeger.Span)
			ctx := metadata.AppendToOutgoingContext(r.Context(), "x-trace-id", fmt.Sprint(jaegerSpan.SpanContext().TraceID()))
			ctx = metadata.AppendToOutgoingContext(ctx, "x-span-id", fmt.Sprint(jaegerSpan.SpanContext().TraceID()))

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			w.Header().Set("x-trace-id", fmt.Sprint(jaegerSpan.SpanContext().TraceID()))

			authToken := extractAuthorizationToken(r.Header)

			if authToken == "" {
				rw := &CapturingResponseWriter{
					ResponseWriter: w,
				}

				r = r.WithContext(ctx)

				next.ServeHTTP(rw, r)

				metrics.ResponseTime.WithLabelValues(*body.OperationName).Observe(float64(time.Since(start).Milliseconds()))
				return

			}

			authResp, err := authClient.ValidateToken(ctx, &auth_service.ValidateTokenRequest{
				Token: authToken,
			})
			if err != nil {

				next.ServeHTTP(w, r)

				metrics.ResponseTime.WithLabelValues(*body.OperationName).Observe(float64(time.Since(start).Milliseconds()))
				return
			}

			session := Session{
				Username: authResp.GetUsername(),
				ClientID: authResp.GetClientId(),
			}
			ctx = context.WithValue(ctx, SessionCtxName, session)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			metrics.ResponseTime.WithLabelValues(*body.OperationName).Observe(float64(time.Since(start).Milliseconds()))
		})
	}
}

func extractAuthorizationToken(headers http.Header) string {
	authToken := strings.Split(
		headers.Get("Authorization"), " ",
	)

	if len(authToken) == 2 {
		return authToken[1]
	}

	cookies := headers.Get("Cookie")
	for _, cookie := range strings.Split(cookies, ";") {
		parts := strings.Split(cookie, "=")
		if len(parts) < 2 {
			continue
		}

		if strings.ToLower(strings.TrimSpace(parts[0])) == "authorization" {
			return parts[1]
		}
	}

	return ""
}
