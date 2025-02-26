package config

import (
	"cdn-service/internal/controller"
	"cdn-service/internal/middleware"
	"cdn-service/internal/services"
	"cdn-service/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"path/filepath"
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
		log.Info().Msg("🛑 Shutting down server...")

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
	go server.initNats(cfg.Nats)
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
	log.Info().Msg("🚀 Starting server...")
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

// ImageDeleteRequest represents the payload for deleting images
type ImageDeleteRequest struct {
	ClientID string   `json:"client_id"` // Client who owns the images
	Images   []string `json:"images"`    // List of images to delete
}

// ImageDeleteResponse represents the response after deletion
type ImageDeleteResponse struct {
	ClientID string   `json:"client_id"`
	Deleted  []string `json:"deleted"`
	Failed   []string `json:"failed"`
}

func (s *ServerConfig) initNats(apiNats string) {
	nc, err := nats.Connect(apiNats)
	if err != nil {
		fmt.Println("❌ Failed to connect to NATS", err)
		return
	}
	defer nc.Close()

	_, err = nc.Subscribe("asset.image.delete", func(msg *nats.Msg) {
		var request ImageDeleteRequest
		err := json.Unmarshal(msg.Data, &request)
		if err != nil {
			fmt.Println("❌ Failed to decode delete request", err)
			return
		}

		// Delete images
		deleted, failed := DeleteImages(request.ClientID, request.Images)
		log.Log().Msgf("Deleted: %v, Failed: %v", deleted, failed)

	})

	select {} // Keep the subscriber running indefinitely
}

// DeleteImages removes files from storage
func DeleteImages(clientID string, images []string) ([]string, []string) {
	var deleted []string
	var failed []string

	for _, img := range images {
		filePath := filepath.Join(utils.UploadDir, clientID, img)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Println("⚠️ File not found:", filePath)
			failed = append(failed, img)
			continue
		}

		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("❌ Failed to delete:", filePath)
			failed = append(failed, img)
		} else {
			fmt.Println("✅ Deleted:", filePath)
			deleted = append(deleted, img)
		}
	}

	return deleted, failed
}
