package config

import (
	"cdn-service/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Config holds application-wide configurations
type Config struct {
	AppPort   string `envconfig:"APP_PORT" default:"8181"`
	JWTSecret string `envconfig:"JWT_SECRET" default:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"`
	RedisHost string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort string `envconfig:"REDIS_PORT" default:"6379"`
	RedisDB   int    `envconfig:"REDIS_DB" default:"0"`
	RedisPass string `envconfig:"REDIS_PASSWORD" default:""`
	Nats      string `envconfig:"NATS_URL" default:"nats://localhost:4222"`
}

func initDir() {
	if _, err := os.Stat(utils.UploadDir); os.IsNotExist(err) {
		err := os.Mkdir(utils.UploadDir, os.ModePerm)
		if err != nil {
			return
		}
	}
}

// LoadConfig loads environment variables into the Config struct
func LoadConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"AppPort":   cfg.AppPort,
		"RedisHost": cfg.RedisHost,
	}).Info("✅ Configuration loaded successfully")

	return &cfg
}

// InitRedis initializes and returns a Redis client with retry logic
func InitRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	// Connection Retry Logic
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			logrus.Info("✅ Connected to Redis")
			return rdb
		}

		logrus.WithFields(logrus.Fields{
			"attempt": i,
			"error":   err.Error(),
		}).Warn("⏳ Retrying Redis connection...")

		time.Sleep(2 * time.Second)
	}

	logrus.Fatal("❌ Failed to connect to Redis after retries")
	return nil
}

// CloseRedis closes the Redis connection properly
func CloseRedis(rdb *redis.Client) {
	if err := rdb.Close(); err != nil {
		logrus.WithError(err).Error("Error closing Redis connection")
	} else {
		logrus.Info("✅ Redis connection closed")
	}
}

// InitGin initializes the Gin engine with appropriate configurations
func InitGin() *gin.Engine {
	// Set Gin mode based on environment
	if ginMode := gin.Mode(); ginMode != gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		logrus.Warn("⚠ Running in DEBUG mode. Use `GIN_MODE=release` in production.")
	} else {
		logrus.Info("✅ Running in RELEASE mode.")
	}

	// Create a new Gin router
	engine := gin.New()

	// Middleware
	engine.Use(gin.Recovery()) // Handles panics and prevents crashes
	engine.Use(gin.Logger())   // Logs HTTP requests

	// Security Headers (Prevents Clickjacking & XSS Attacks)
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Next()
	})

	logrus.Info("🚀 Gin HTTP server initialized successfully")
	return engine
}
