package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmehra2102/StudySync/internal/config"
	httpServer "github.com/dmehra2102/StudySync/internal/delivery/http"
	"github.com/dmehra2102/StudySync/internal/repository"
	"github.com/dmehra2102/StudySync/internal/service"
	"github.com/dmehra2102/StudySync/pkg/database"
	"github.com/dmehra2102/StudySync/pkg/logger"
	"github.com/dmehra2102/StudySync/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log := logger.New(cfg.Log.Level)

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	redisClient := redis.New(cfg.Redis)
	defer redisClient.Close()

	// Initialize repositories
	sessionRepo := repository.NewStudySessionRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Initialize service
	notificationService := service.NewNotificationService(notificationRepo, redisClient)

	schedulerService := service.NewSchedulerService(sessionRepo, taskRepo, notificationRepo, notificationService)

	schedulerService.Start()
	defer schedulerService.Stop()

	// Initialize HTTP server
	server := httpServer.NewServer(cfg, db, redisClient, log)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	log.Info().Msgf("Server started on port %s", cfg.App.Port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
