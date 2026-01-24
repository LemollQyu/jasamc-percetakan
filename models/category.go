package models

import (
	"time"

	"gorm.io/datatypes"
)

// model kategori jasa
type CategoryJasa struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Slug        string         `json:"slug"`
	IsActive    bool           `json:"is_active"`
	Meta        datatypes.JSON `json:"meta"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTimeUTC"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTimeUTC"`
}

// param untuk create category jasa
type ParamCreateCategoryJasa struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Slug        string         `json:"slug"`
	Meta        datatypes.JSON `json:"meta"`
}

type ParamUpdateCategoryJasa struct {
	ID          int64  `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	IsActive    *bool  `json:"is_active"`
}

// param untuk Delete category
type ParamDeleteCategoryJasa struct {
	ID int64 `json:"id" binding:"required"`
}
