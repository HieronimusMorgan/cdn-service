package controller

import (
	"cdn-service/internal/services"
	"cdn-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type ImageController interface {
	UploadImages(context *gin.Context)
	GetImage(context *gin.Context)
}

type imageController struct {
	ImageService services.ImageService
	JWTService   utils.JWTService
}

func NewImageController(imageService services.ImageService, jwtService utils.JWTService) ImageController {
	return imageController{ImageService: imageService, JWTService: jwtService}
}

func (h imageController) UploadImages(context *gin.Context) {
	// Parse form data
	err := context.Request.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse form")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Extract JWT Token
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Get multiple files from the form
	form, _ := context.MultipartForm()
	files := form.File["images"] // Get all uploaded images
	log.Info().Msgf("Uploading %d images", len(files))
	// If no images are uploaded
	if len(files) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	// Call service to upload images
	uploadedImages, err := h.ImageService.UploadImages(files, token.ClientID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return uploaded images response
	context.JSON(http.StatusCreated, gin.H{"data": uploadedImages})
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
