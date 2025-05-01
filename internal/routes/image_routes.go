package routes

import (
	"cdn-service/config"
	"cdn-service/internal/controller"
	"github.com/gin-gonic/gin"
)

func ImageRoutes(r *gin.Engine, middleware config.Middleware, controller controller.ImageController) {

	routerGroup := r.Group("/v1")
	{
		routerGroup.GET("/cdn/:clientID/:filename", controller.GetImage)
	}
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.POST("/upload-photo-profile", controller.UploadPhotoProfile)
		routerGroup.POST("/upload", controller.UploadImages)
	}
}
