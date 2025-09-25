package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsMiddleware возвращает middleware для сбора RED-метрик
func MetricsMiddleware() func(http.Handler) http.Handler {

	// Регистрируем метрики
	var (
		httpRequestsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests by method, path and status code.",
			},
			[]string{"method", "path", "status_code"},
		)

		httpRequestDuration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: []float64{0.1, 0.3, 0.5, 1, 2, 5}, // Настраиваемые бакеты
			},
			[]string{"method", "path", "status_code"},
		)

		httpRequestsInProgress = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_requests_in_progress",
				Help: "Current number of HTTP requests being processed.",
			},
			[]string{"method", "path"},
		)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("call metrics middleware")

			// Засекаем время начала обработки
			start := time.Now()

			// Получаем путь шаблона маршрута (например, "/users/{id}")
			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern == "" {
				routePattern = "unknown"
			}

			// Увеличиваем счетчик текущих запросов
			httpRequestsInProgress.WithLabelValues(r.Method, routePattern).Inc()
			// Уменьшаем при завершении
			defer httpRequestsInProgress.WithLabelValues(r.Method, routePattern).Dec()

			// Создаем кастомный ResponseWriter для захвата статус кода
			ww := &responseWriter{w: w, statusCode: http.StatusOK}

			// Обрабатываем запрос
			next.ServeHTTP(ww, r)

			// Вычисляем длительность
			duration := time.Since(start).Seconds()

			// Преобразуем статус код в строку
			statusCode := strconv.Itoa(ww.statusCode)

			// Записываем метрики
			httpRequestsTotal.WithLabelValues(r.Method, routePattern, statusCode).Inc()
			httpRequestDuration.WithLabelValues(r.Method, routePattern, statusCode).Observe(duration)
		})
	}
}

// responseWriter оборачивает http.ResponseWriter для захвата статус кода
type responseWriter struct {
	w          http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.w.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.w.Write(b)
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.w
}
