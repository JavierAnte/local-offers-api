package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/config"
	"github.com/JavierAnte/local-offers-api/internal/database"
	"github.com/JavierAnte/local-offers-api/internal/handlers"
	"github.com/JavierAnte/local-offers-api/internal/repositories"
	"github.com/JavierAnte/local-offers-api/internal/server"
	"github.com/JavierAnte/local-offers-api/internal/services"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("access database connection: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("close database connection: %v", err)
		}
	}()

	const uploadDir = "uploads"

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	offerRepo := repositories.NewOfferRepository(db)
	offerService := services.NewOfferService(offerRepo)
	offerHandler := handlers.NewOfferHandler(offerService)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)
	commentHandler := handlers.NewCommentHandler(commentService)

	offerVoteRepo := repositories.NewOfferVoteRepository(db)
	offerVoteService := services.NewOfferVoteService(offerVoteRepo)
	offerVoteHandler := handlers.NewOfferVoteHandler(offerVoteService)

	uploadService := services.NewUploadService(uploadDir)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	r := server.NewRouter(cfg.JWTSecret, uploadDir, server.Handlers{
		Auth:      authHandler,
		Offer:     offerHandler,
		Comment:   commentHandler,
		OfferVote: offerVoteHandler,
		Upload:    uploadHandler,
	})

	httpServer := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server running on %s", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	case <-shutdownSignal.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownContext); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
			if err := httpServer.Close(); err != nil {
				log.Printf("forced server close failed: %v", err)
			}
		}
	}
}
