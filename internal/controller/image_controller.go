package controller

import (
	"cdn-service/internal/services"
	"cdn-service/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ImageController interface {
	UploadImage(context *gin.Context)
	GetImage(context *gin.Context)
}

type imageController struct {
	ImageService services.ImageService
	JWTService   utils.JWTService
}

func NewImageController(imageService services.ImageService, jwtService utils.JWTService) ImageController {
	return imageController{ImageService: imageService, JWTService: jwtService}
}

func (h imageController) UploadImage(context *gin.Context) {
	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	image, err := h.ImageService.UploadImage(file, token.ClientID)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(201, gin.H{"data": image})
}

func (h imageController) GetImage(context *gin.Context) {
	clientID := context.Param("clientID")
	filename := context.Param("filename")
	if clientID == "" || filename == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID or filename"})
		return
	}
	image, err := h.ImageService.GetImage(filename, clientID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contentType := "application/octet-stream"
	if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filename, ".webp") {
		contentType = "image/webp"
	}

	// Set cache headers for faster performance
	context.Header("Cache-Control", "public, max-age=86400") // Cache for 1 day
	context.Header("Content-Type", contentType)

	// Serve the file
	context.File(*image)
}
