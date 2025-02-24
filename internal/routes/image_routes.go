package routes

import (
	"cdn-service/config"
	"cdn-service/internal/controller"
	"github.com/gin-gonic/gin"
)

func ImageRoutes(r *gin.Engine, middleware config.Middleware, controller controller.ImageController) {

	routerGroup := r.Group("/v1")
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.POST("/upload", controller.UploadImage)
		routerGroup.GET("/cdn/:clientID/:filename", controller.GetImage)
	}
}
