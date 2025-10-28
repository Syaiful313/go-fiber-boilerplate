package services

import (
	"errors"
	"log"
	"mime/multipart"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/pkg/pagination"

	"gorm.io/gorm"
)

type SampleService struct {
	cfg               *config.Config
	cloudinaryService *CloudinaryService
}

func NewSampleService(cfg *config.Config) *SampleService {
	cloudinaryService, err := NewCloudinaryService(cfg)
	if err != nil {
		log.Printf("Cloudinary service disabled: %v", err)
	}
	return &SampleService{
		cfg:               cfg,
		cloudinaryService: cloudinaryService,
	}
}

func (s *SampleService) GetSamples(params pagination.Params) ([]models.Sample, pagination.Meta, error) {
	var samples []models.Sample
	var total int64

	if err := database.GetDB().Model(&models.Sample{}).Count(&total).Error; err != nil {
		return nil, pagination.Meta{}, err
	}

	if params.All {
		if err := database.GetDB().
			Preload("User").
			Order(params.OrderClause("created_at")).
			Find(&samples).Error; err != nil {
			return nil, pagination.Meta{}, err
		}

		meta := pagination.Meta{
			HasNext:     false,
			HasPrevious: false,
			Page:        1,
			PerPage:     int(total),
			Total:       total,
		}
		return samples, meta, nil
	}

	if err := database.GetDB().
		Preload("User").
		Offset(params.Offset()).
		Limit(params.PerPage).
		Order(params.OrderClause("created_at")).
		Find(&samples).Error; err != nil {
		return nil, pagination.Meta{}, err
	}

	meta := pagination.BuildMeta(total, params)
	return samples, meta, nil
}

func (s *SampleService) GetSampleById(id int) (*models.Sample, error) {
	var sample models.Sample
	if err := database.GetDB().
		Preload("User").
		First(&sample, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &sample, nil
}

func (s *SampleService) CreateSample(userID uint, req models.CreateSampleRequest, imageFile *multipart.FileHeader) (*models.Sample, error) {
	var existingBlog models.Sample
	if err := database.GetDB().Where("title = ?", req.Title).First(&existingBlog).Error; err == nil {
		return nil, errors.New("title already exists")
	}
	sample := models.Sample{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}

	// Handle image upload if provided
	if imageFile != nil && s.cloudinaryService != nil {
		uploadResult, err := s.cloudinaryService.UploadImage(imageFile, "samples")
		if err != nil {
			return nil, errors.New("failed to upload image: " + err.Error())
		}
		sample.ImageURL = uploadResult.SecureURL
		sample.ImagePublicID = uploadResult.PublicID
	}

	if err := database.GetDB().Create(&sample).Error; err != nil {
		// Cleanup image if database save fails
		if sample.ImagePublicID != "" && s.cloudinaryService != nil {
			s.cloudinaryService.DeleteImage(sample.ImagePublicID)
		}
		return nil, err
	}

	database.GetDB().Preload("User").First(&sample, sample.ID)
	return &sample, nil
}

func (s *SampleService) UpdateSample(userID uint, id int, req models.UpdateSampleRequest, imageFile *multipart.FileHeader) (*models.Sample, error) {
	var sample models.Sample
	if err := database.GetDB().First(&sample, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if sample.UserID != userID {
		return nil, errors.New("forbidden")
	}

	oldImagePublicID := sample.ImagePublicID

	if req.Title != "" {
		sample.Title = req.Title
	}
	if req.Description != "" {
		sample.Description = req.Description
	}

	// Handle new image upload
	if imageFile != nil && s.cloudinaryService != nil {
		uploadResult, err := s.cloudinaryService.UploadImage(imageFile, "samples")
		if err != nil {
			return nil, errors.New("failed to upload image: " + err.Error())
		}
		sample.ImageURL = uploadResult.SecureURL
		sample.ImagePublicID = uploadResult.PublicID
	}

	if err := database.GetDB().Save(&sample).Error; err != nil {
		return nil, err
	}

	// Delete old image if new one was uploaded successfully
	if imageFile != nil && oldImagePublicID != "" && s.cloudinaryService != nil {
		s.cloudinaryService.DeleteImage(oldImagePublicID)
	}

	database.GetDB().Preload("User").First(&sample, sample.ID)
	return &sample, nil
}

func (s *SampleService) DeleteSample(userID uint, id int) error {
	var sample models.Sample
	if err := database.GetDB().First(&sample, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	if sample.UserID != userID {
		return errors.New("forbidden")
	}

	if err := database.GetDB().Delete(&sample).Error; err != nil {
		return err
	}

	return nil
}
