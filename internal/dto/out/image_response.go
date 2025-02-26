package out

import "time"

// ImageResponse represents the response data for uploaded images
type ImageResponse struct {
	ImageURL   string    `json:"image_url"`   // URL or file path of the uploaded image
	FileType   string    `json:"file_type"`   // Image format (jpg, png, etc.)
	FileSize   int64     `json:"file_size"`   // Image size in bytes
	UploadedAt time.Time `json:"uploaded_at"` // Timestamp of upload
}
