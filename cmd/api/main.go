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
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	_ = db

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	offerRepo := repositories.NewOfferRepository(db)
	offerService := services.NewOfferService(offerRepo)
	offerHandler := handlers.NewOfferHandler(offerService)

	r.Post("/offers", offerHandler.Create)
	r.Get("/offers/{id}", offerHandler.FindByID)
	r.Get("/offers/nearby", offerHandler.FindNearby)

	fmt.Printf("Server running on :%s\n", cfg.AppPort)

	http.ListenAndServe(":"+cfg.AppPort, r)
}
