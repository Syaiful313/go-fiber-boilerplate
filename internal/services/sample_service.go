package services

import (
	"errors"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/models"

	"gorm.io/gorm"
)

type SampleService struct{
	cfg *config.Config
}

func NewSampleService(cfg *config.Config) *SampleService {
	return &SampleService{cfg: cfg}
}

func (s *SampleService) GetSamples(page, limit int) ([]models.Sample, int64, error) {
	var samples []models.Sample
	offset := (page - 1) * limit

	if err := database.GetDB().
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&samples).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := database.GetDB().Model(&models.Sample{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return samples, total, nil
}

func (s *SampleService) GetSample(id int) (*models.Sample, error) {
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

func (s *SampleService) CreateSample(userID uint, req models.CreateSampleRequest) (*models.Sample, error) {
	sample := models.Sample{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}

	if err := database.GetDB().Create(&sample).Error; err != nil {
		return nil, err
	}

	database.GetDB().Preload("User").First(&sample, sample.ID)
	return &sample, nil
}

func (s *SampleService) UpdateSample(userID uint, id int, req models.UpdateSampleRequest) (*models.Sample, error) {
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

	if req.Title != "" {
		sample.Title = req.Title
	}
	if req.Description != "" {
		sample.Description = req.Description
	}

	if err := database.GetDB().Save(&sample).Error; err != nil {
		return nil, err
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
