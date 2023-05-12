package gsmiddleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-chi/chi/v5"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_with_status",
			Help: "Number of requests resulting in each status code",
		},
		[]string{"path", "status"},
	)
	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_response_time_seconds",
			Help: "Duration of HTTP requests made to this API",
		},
		[]string{"path"},
	)
)

func MetricsHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ww := NewResponseWriterWrapper(w, false)
		start := time.Now()

		// Call the next handler in Handler chain.
		next.ServeHTTP(ww, r)

		// Get the path template from the RouteContext.
		ctx := chi.RouteContext(r.Context())
		if ctx != nil {
			pathTemplate := ctx.RoutePattern()
			if pathTemplate != "" {
				// Get response status code and increment metric.
				code := ww.statusCode
				totalRequests.WithLabelValues(pathTemplate, strconv.Itoa(*code)).Inc()
				responseTime.WithLabelValues(pathTemplate).Observe(float64(time.Since(start).Seconds()))
			}
		}
	}
	return http.HandlerFunc(fn)
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseTime)
}