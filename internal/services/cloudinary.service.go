package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"go-fiber-boilerplate/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	cld *cloudinary.Cloudinary
	cfg *config.Config
}

type UploadResult struct {
	PublicID  string `json:"public_id"`
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Bytes     int    `json:"bytes"`
}

func NewCloudinaryService(cfg *config.Config) (*Cloudinary, error) {
	if cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == "" {
		return nil, fmt.Errorf("cloudinary credentials are required")
	}

	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &Cloudinary{
		cld: cld,
		cfg: cfg,
	}, nil
}

func (s *Cloudinary) UploadImage(file *multipart.FileHeader, folder string) (*UploadResult, error) {
	if !s.isImageFile(file.Filename) {
		return nil, fmt.Errorf("file must be an image (jpg, jpeg, png, gif, webp)")
	}

	maxSize := int64(5 * 1024 * 1024)
	if file.Size > maxSize {
		return nil, fmt.Errorf("file size exceeds 5MB limit")
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	timestamp := time.Now().Unix()
	filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	publicID := fmt.Sprintf("%s_%d", filename, timestamp)

	overwrite := true
	uploadParams := uploader.UploadParams{
		PublicID:       publicID,
		Folder:         folder,
		Transformation: "f_auto,q_auto,w_1200",
		Overwrite:      &overwrite,
	}

	result, err := s.cld.Upload.Upload(context.Background(), src, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return &UploadResult{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     result.Bytes,
	}, nil
}

func (s *Cloudinary) DeleteImage(publicID string) error {
	_, err := s.cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from cloudinary: %w", err)
	}
	return nil
}

func (s *Cloudinary) GenerateURL(publicID string, size string) string {
	transformations := map[string]string{
		"thumbnail": "w_150,h_150,c_thumb,f_auto,q_auto",
		"small":     "w_300,f_auto,q_auto",
		"medium":    "w_600,f_auto,q_auto",
		"large":     "w_1200,f_auto,q_auto",
		"original":  "f_auto,q_auto",
	}

	transform, exists := transformations[size]
	if !exists {
		transform = transformations["original"]
	}

	baseURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload", s.cfg.CloudinaryCloudName)
	return fmt.Sprintf("%s/%s/%s", baseURL, transform, publicID)
}

func (s *Cloudinary) GetImageVariants(publicID string) map[string]string {
	sizes := []string{"thumbnail", "small", "medium", "large", "original"}
	variants := make(map[string]string)

	for _, size := range sizes {
		variants[size] = s.GenerateURL(publicID, size)
	}

	return variants
}

func (s *Cloudinary) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
