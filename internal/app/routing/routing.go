package routing

import (
	"ms_dialog/internal/app/handlers"
	"ms_dialog/internal/app/middleware/metrics"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouting(handlers *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// объявляем MiddleWare для сбора метрик
	r.Use(metrics.MetricsMiddleware())

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/v2/dialog/{user_id}/send", handlers.SendMsgUser)
	r.Get("/v2/dialog/{user_id}/list/{error}", handlers.GetDialog)

	// Health check и метрики
	r.Get("/health", handlers.HealthHandler)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/v2/test", handlers.Test)

	return r
}
