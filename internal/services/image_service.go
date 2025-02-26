package services

import (
	response "cdn-service/internal/dto/out"
	"cdn-service/internal/utils"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageService defines the interface for managing asset categories
type ImageService interface {
	UploadImages(files []*multipart.FileHeader, clientID string) ([]response.ImageResponse, error)
	GetImage(imageUrl, clientID string) (*string, error)
}

// imageService implements ImageService
type imageService struct {
	Redis utils.RedisService
}

// NewImageService initializes the service
func NewImageService(redis utils.RedisService) ImageService {
	return &imageService{
		Redis: redis,
	}
}

func (s *imageService) UploadImages(files []*multipart.FileHeader, clientID string) ([]response.ImageResponse, error) {
	var uploadedImages []response.ImageResponse
	uploadDir := filepath.Join("./uploads", clientID)

	log.Info().Msgf("Uploading %d images for client: %s", len(files), clientID)
	// Ensure client directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create directory")
			return nil, err
		}
	}

	// Process each file
	for _, file := range files {
		// Open file
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to close file")
				return
			}
		}(src)

		// Generate safe filename
		extension := strings.ToLower(filepath.Ext(file.Filename))
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		filePath := filepath.Join(uploadDir, newFileName)

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to close file")
				return
			}
		}(dst)

		// Copy file data
		fileSize, err := io.Copy(dst, src)
		if err != nil {
			return nil, err
		}
		log.Info().Msgf("File uploaded: %s (%d bytes)", newFileName, fileSize)

		// Get file metadata
		imageResponse := response.ImageResponse{
			ImageURL:   fmt.Sprintf("/cdn/%s/%s", clientID, newFileName), // Serve via API
			FileType:   strings.TrimPrefix(extension, "."),               // Remove dot (jpg, png)
			FileSize:   fileSize,
			UploadedAt: time.Now(),
		}

		// Append to response
		uploadedImages = append(uploadedImages, imageResponse)
	}

	return uploadedImages, nil
}

func (s *imageService) GetImage(imageUrl, clientID string) (*string, error) {

	filePath := filepath.Join(utils.UploadDir, clientID, imageUrl)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("file not found")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("file not found")
	}
	return &filePath, nil
}
