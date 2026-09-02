package main

import (
	"fmt"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/config"
	"github.com/JavierAnte/local-offers-api/internal/database"
	"github.com/JavierAnte/local-offers-api/internal/handlers"
	"github.com/JavierAnte/local-offers-api/internal/repositories"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// Mobile clients (Expo Go, native builds) don't send a browser Origin,
		// and there's no cookie-based auth yet, so a permissive origin is safe here.
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	offerRepo := repositories.NewOfferRepository(db)
	offerService := services.NewOfferService(offerRepo)
	offerHandler := handlers.NewOfferHandler(offerService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/offers", offerHandler.Create)
		r.Get("/offers/{id}", offerHandler.FindByID)
		r.Get("/offers/nearby", offerHandler.FindNearby)
	})

	fmt.Printf("Server running on :%s\n", cfg.AppPort)

	http.ListenAndServe(":"+cfg.AppPort, r)
}
