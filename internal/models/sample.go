package models

import (
	"time"

	"gorm.io/gorm"
)

type Sample struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	UserID      uint           `json:"userId" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type SampleResponse struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	UserID      uint         `json:"userId"`
	User        UserResponse `json:"user"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type CreateSampleRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateSampleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (s *Sample) ToResponse() SampleResponse {
	return SampleResponse{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		UserID:      s.UserID,
		User:        s.User.ToResponse(),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

