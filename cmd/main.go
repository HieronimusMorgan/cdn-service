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
	//assets.AssetCategoryRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetCategory)
	//assets.AssetStatusRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetStatus)
	//assets.AssetRoutes(engine, serverConfig.Middleware, serverConfig.Controller.Asset)
	//assets.AssetWishlistRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetWishlist)
	//assets.AssetMaintenanceRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetMaintenance)
	//assets.AssetMaintenanceTypeRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetMaintenanceType)
	// Run server
	log.Println("Starting server on :8181")
	err = engine.Run(":8181")
	if err != nil {
		return
	}
}
