package models

import (
	"time"

	"gorm.io/datatypes"
)

type FullService struct {
	ID            int64                  `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID    int64                  `json:"category_id" binding:"required" gorm:"not null"`
	Name          string                 `json:"name" binding:"required" gorm:"type:varchar(200);not null"`
	Slug          string                 `json:"slug" gorm:"type:varchar(200);not null"`
	Description   string                 `json:"description" binding:"required" gorm:"type:text;not null"`
	BasePrice     float64                `json:"base_price" gorm:"type:numeric(12,2)"`
	IsActive      bool                   `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time              `json:"created_at" gorm:"autoCreateTimeUTC"`
	UpdatedAt     time.Time              `json:"updated_at" gorm:"autoUpdateTimeUTC"`
	Category      Category               `json:"category"`
	Media         []ServiceMedia         `json:"media" gorm:"foreignKey:ServiceID"` // slice untuk media
	Spesification []ServiceSpesification `json:"spesification" gorm:"foreignKey:ServiceID"`
}

type Category struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Slug        string         `json:"slug"`
	IsActive    bool           `json:"is_active"`
	Meta        datatypes.JSON `json:"meta"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (FullService) TableName() string {
	return "services"
}

func (Category) TableName() string {
	return "service_categories"
}
