package http

import (
	"context"
	"net/http"
	"time"

	"github.com/dmehra2102/StudySync/internal/config"
	"github.com/dmehra2102/StudySync/internal/delivery/http/handlers"
	"github.com/dmehra2102/StudySync/internal/delivery/http/middleware"
	"github.com/dmehra2102/StudySync/internal/repository"
	"github.com/dmehra2102/StudySync/internal/service"
	"github.com/dmehra2102/StudySync/pkg/auth"
	"github.com/dmehra2102/StudySync/pkg/logger"
	"github.com/dmehra2102/StudySync/pkg/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	cfg        *config.Config
	db         *gorm.DB
	redis      *redis.Client
	log        *logger.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config, db *gorm.DB, redis *redis.Client, log *logger.Logger) *Server {
	server := &Server{
		cfg:   cfg,
		db:    db,
		redis: redis,
		log:   log,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	if s.cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	s.httpServer = &http.Server{
		Addr:    ":" + s.cfg.App.Port,
		Handler: router,
	}

	// middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(s.log))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestTimeout(5 * time.Second))

	// Initialize dependencies
	jwtAuth := auth.NewJWTAuth(s.cfg.Auth.JWTSecret, s.cfg.Auth.JWTExpiry)
	authMiddleware := middleware.AuthMiddleware(jwtAuth)

	// Repositories
	userRepo := repository.NewUserRepository(s.db)
	sessionRepo := repository.NewStudySessionRepository(s.db)
	taskRepo := repository.NewTaskRepository(s.db)
	notificationRepo := repository.NewNotificationRepository(s.db)

	// Services
	userService := service.NewUserService(userRepo, jwtAuth)
	sessionService := service.NewStudySessionService(sessionRepo)
	taskService := service.NewTaskService(taskRepo)
	notificationService := service.NewNotificationService(notificationRepo, s.redis)

	// Handlers
	authHandler := handlers.NewUserHandler(userService)
	sessionHandler := handlers.NewStudySessionHandler(sessionService)
	taskHandler := handlers.NewTaskHandler(taskService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Public Routes
	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}
	}

	// Protected Routes
	authGroup := v1.Group("")
	authGroup.Use(authMiddleware)
	{
		userGroup := authGroup.Group("/users")
		{
			userGroup.GET("/profile", authHandler.GetProfile)
		}

		sessionGroup := authGroup.Group("/study-sessions")
		{
			sessionGroup.POST("", sessionHandler.CreateSession)
			sessionGroup.GET("", sessionHandler.GetSessions)
			sessionGroup.GET("/upcoming", sessionHandler.GetUpcomingSessions)
			sessionGroup.GET("/stats", sessionHandler.GetStats)
			sessionGroup.GET("/:id", sessionHandler.GetSession)
			sessionGroup.PUT("/:id", sessionHandler.UpdateSession)
			sessionGroup.DELETE("/:id", sessionHandler.DeleteSession)
		}

		taskGroup := authGroup.Group("/tasks")
		{
			taskGroup.POST("", taskHandler.CreateTask)
			taskGroup.GET("", taskHandler.GetTasks)
			taskGroup.GET("/overdue", taskHandler.GetOverdueTasks)
			taskGroup.GET("/:id", taskHandler.GetTask)
			taskGroup.PUT("/:id", taskHandler.UpdateTask)
			taskGroup.DELETE("/:id", taskHandler.DeleteTask)
		}

		notificationGroup := authGroup.Group("/notifications")
		{
			notificationGroup.GET("", notificationHandler.GetNotifications)
			notificationGroup.PUT("/:id/read", notificationHandler.MarkAsRead)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "OK",
			"environment": s.cfg.App.Env,
		})
	})

}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
