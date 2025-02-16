package services

import (
	"cdn-service/internal/utils"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
)

// ImageService defines the interface for managing asset categories
type ImageService interface {
	UploadImage(file *multipart.FileHeader, clientID string) (interface{}, error)
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

// UploadImage uploads an image for a specific client
func (s *imageService) UploadImage(file *multipart.FileHeader, clientID string) (interface{}, error) {
	if clientID == "" {
		return nil, errors.New("clientID is required")
	}

	clientDir := filepath.Join(utils.UploadDir, clientID)

	// Ensure client directory exists
	if err := os.MkdirAll(clientDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := src.Close(); err != nil {
			// Log the error instead of ignoring it
			// log.Println("Error closing src file:", err)
		}
	}()

	// Decode image
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	// Resize Image
	dst := image.NewRGBA(image.Rect(0, 0, 800, 600)) // Resize to 800x600
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// Generate unique filename
	newFileName := uuid.New().String() + "." + format
	filePath := filepath.Join(clientDir, newFileName)

	// Create file to save the optimized image
	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := out.Close(); err != nil {
			// log.Println("Error closing output file:", err)
		}
	}()

	// Encode and save image
	switch format {
	case "png":
		if err := png.Encode(out, dst); err != nil {
			return nil, err
		}
	case "jpeg", "jpg":
		if err := jpeg.Encode(out, dst, &jpeg.Options{Quality: 80}); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported image format: " + format)
	}

	// Return URL path
	return "/cdn/" + clientID + "/" + newFileName, nil
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
