package config

import (
	"cdn-service/internal/controller"
	"cdn-service/internal/middleware"
	"cdn-service/internal/services"
	"cdn-service/internal/utils"
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewServerConfig() (*ServerConfig, error) {
	initDir()
	cfg := LoadConfig()
	redisClient := InitRedis(cfg)
	redisService := utils.NewRedisService(*redisClient)
	engine := InitGin()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down gracefully...")

		// Close database and Redis before exiting
		CloseRedis(redisClient)

		os.Exit(0)
	}()

	server := &ServerConfig{
		Gin:        engine,
		Config:     cfg,
		Redis:      redisService,
		JWTService: utils.NewJWTService(cfg.JWTSecret),
	}

	server.initRepository()
	server.initServices()
	server.initController()
	server.initMiddleware()
	return server, nil
}

// initRepository initializes database access objects (Repository)
func (s *ServerConfig) initRepository() {
	s.Repository = Repository{}
}

// initServices initializes the application services
func (s *ServerConfig) initServices() {
	s.Services = Services{
		ImageService: services.NewImageService(s.Redis),
	}
}

// Start initializes everything and returns an error if something fails
func (s *ServerConfig) Start() error {
	log.Println("✅ Server configuration initialized successfully!")
	return nil
}

func (s *ServerConfig) initController() {
	s.Controller = Controller{
		ImageController: controller.NewImageController(s.Services.ImageService, s.JWTService),
	}
}

func (s *ServerConfig) initMiddleware() {
	s.Middleware = Middleware{
		AuthMiddleware: middleware.NewAuthMiddleware(s.JWTService),
	}
}
