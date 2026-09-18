package server

import (
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handlers struct {
	Auth      *handlers.AuthHandler
	Offer     *handlers.OfferHandler
	Comment   *handlers.CommentHandler
	OfferVote *handlers.OfferVoteHandler
	Upload    *handlers.UploadHandler
}

func NewRouter(jwtSecret string, uploadDir string, h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", h.Auth.Register)
		r.Post("/auth/login", h.Auth.Login)

		r.Get("/offers/{id}", h.Offer.FindByID)
		r.Get("/offers/nearby", h.Offer.FindNearby)
		r.Get("/offers/{id}/comments", h.Comment.List)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(jwtSecret))
			r.Post("/offers", h.Offer.Create)
			r.Post("/offers/{id}/comments", h.Comment.Create)
			r.Post("/offers/{id}/votes", h.OfferVote.Create)
			r.Post("/uploads", h.Upload.Create)
		})
	})

	return r
}
