package main

import (
	"cdn-service/config"
	"cdn-service/internal/routes"
	"log"
)

func main() {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	if err := serverConfig.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}

	engine := serverConfig.Gin

	routes.ImageRoutes(engine, serverConfig.Middleware, serverConfig.Controller.ImageController)

	// Run server
	log.Println("Starting server on :8181")
	err = engine.Run(":8181")
	if err != nil {
		return
	}
}
