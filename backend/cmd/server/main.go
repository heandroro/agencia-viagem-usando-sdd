package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heandroro/agencia-viagem/backend/internal/api/handlers"
	"github.com/heandroro/agencia-viagem/backend/internal/domain/reservation"
	"github.com/heandroro/agencia-viagem/backend/internal/infra/config"
	"github.com/heandroro/agencia-viagem/backend/internal/infra/database"
	"github.com/heandroro/agencia-viagem/backend/pkg/crypto"
)

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Configurar modo do Gin
	if cfg.APM.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Conectar ao PostgreSQL
	dbClient, err := database.NewPostgresClient(database.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Database: cfg.Postgres.Database,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	log.Println("Connected to PostgreSQL")

	// Inicializar serviço de criptografia
	cryptoService, err := crypto.NewAES256Service(cfg.Security.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize crypto service: %v", err)
	}

	// Inicializar repositórios
	reservationRepo := reservation.NewPostgresRepository(dbClient)
	availabilityRepo := reservation.NewPostgresAvailabilityRepository(dbClient)

	// Criar router
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := dbClient.Health(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Reservations endpoints
		reservations := v1.Group("/reservations")
		{
			reservations.POST("", createReservationHandler(reservationRepo, availabilityRepo, cryptoService))
			reservations.PUT("/:id/travelers", updateTravelersHandler(reservationRepo, cryptoService))
			reservations.GET("/:id/summary", getReservationSummaryHandler(reservationRepo))
		}
	}

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// createReservationHandler cria o handler para POST /reservations
func createReservationHandler(
	repo reservation.Repository,
	availRepo reservation.AvailabilityRepository,
	crypto crypto.Service,
) gin.HandlerFunc {
	handler := handlers.NewCreateReservationHandler(repo, availRepo, crypto)
	return handler.Handle
}

func updateTravelersHandler(
	repo reservation.Repository,
	crypto crypto.Service,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
	}
}

func getReservationSummaryHandler(
	repo reservation.Repository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
	}
}
