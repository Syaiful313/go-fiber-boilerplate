package services

import (
	"errors" 

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/pkg/pagination"

	"gorm.io/gorm"
)

type SampleService struct {
	cfg *config.Config
}

func NewSampleService(cfg *config.Config) *SampleService {
	return &SampleService{cfg: cfg}
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
