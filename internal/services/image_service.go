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
	UploadPhotoProfile(files *multipart.FileHeader, clientID string) (response.ImageResponse, error)
	UploadImages(files []*multipart.FileHeader, clientID string) ([]response.ImageResponse, error)
	GetImage(imageUrl, clientID string) (*string, error)
}

// imageService implements ImageService
type imageService struct {
	Redis            utils.RedisService
	UploadDir        string
	ProfileUploadDir string
}

// NewImageService initializes the service
func NewImageService(redis utils.RedisService, uploadDir string, profileUploadDir string) ImageService {
	return &imageService{
		Redis:            redis,
		UploadDir:        uploadDir,
		ProfileUploadDir: profileUploadDir,
	}
}

// UploadPhotoProfile handles the upload of a single photo profile
func (s *imageService) UploadPhotoProfile(file *multipart.FileHeader, clientID string) (response.ImageResponse, error) {
	uploadBaseDir := s.ProfileUploadDir
	if uploadBaseDir == "" {
		uploadBaseDir = "./image/profile"
	}
	uploadDir := filepath.Join(uploadBaseDir, clientID)

	log.Info().Msgf("Uploading photo profile for client: %s", clientID)

	// Ensure client directory exists (create if not)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Error().Err(err).Msg("Failed to create directory")
			return response.ImageResponse{}, err
		}
	} else {
		// Delete all existing files in the folder
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list existing files for deletion")
			return response.ImageResponse{}, err
		}
		for _, f := range files {
			err := os.Remove(filepath.Join(uploadDir, f.Name()))
			if err != nil {
				log.Error().Err(err).Msgf("Failed to delete old file: %s", f.Name())
			}
		}
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return response.ImageResponse{}, err
	}

	err = src.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close file")
		return response.ImageResponse{}, err
	}

	// Generate safe filename
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(uploadDir, newFileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return response.ImageResponse{}, err
	}

	err = src.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close file")
		return response.ImageResponse{}, err
	}

	// Copy file data
	fileSize, err := io.Copy(dst, src)
	if err != nil {
		return response.ImageResponse{}, err
	}
	log.Info().Msgf("File uploaded: %s (%d bytes)", newFileName, fileSize)

	imageResponse := response.ImageResponse{
		ImageURL: fmt.Sprintf("/cdn/%s/%s", clientID, newFileName),
		FileType: strings.TrimPrefix(filepath.Ext(file.Filename), "."),
		FileSize: fileSize,
	}

	return imageResponse, nil
}

func (s *imageService) UploadImages(files []*multipart.FileHeader, clientID string) ([]response.ImageResponse, error) {
	var uploadedImages []response.ImageResponse

	uploadBaseDir := s.UploadDir
	if uploadBaseDir == "" {
		uploadBaseDir = "./image/asset" // fallback default
	}
	uploadDir := filepath.Join(uploadBaseDir, clientID)

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

func (s *imageService) GetImage(imageName, clientID string) (*string, error) {
	// Check in UPLOAD_DIR
	uploadBaseDir := s.UploadDir
	if uploadBaseDir == "" {
		uploadBaseDir = "./uploads"
	}
	filePath := filepath.Join(uploadBaseDir, clientID, imageName)

	if _, err := os.Stat(filePath); err == nil {
		return &filePath, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// If not found, check in PROFILE_UPLOAD_DIR
	profileUploadDir := s.ProfileUploadDir
	if profileUploadDir == "" {
		profileUploadDir = "./profile_uploads"
	}
	profilePath := filepath.Join(profileUploadDir, clientID, imageName)

	if _, err := os.Stat(profilePath); err == nil {
		return &profilePath, nil
	} else if os.IsNotExist(err) {
		return nil, errors.New("file not found")
	} else {
		return nil, err
	}
}
