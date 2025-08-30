package routing

import (
	"ms_dialog/internal/app/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func NewRouting(handlers *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

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
	r.Get("/v2/dialog/{user_id}/list", handlers.GetDialog)

	return r
}
