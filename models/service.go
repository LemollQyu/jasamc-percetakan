package models

import (
	"time"

	"gorm.io/datatypes"
)

type Service struct {
	ID            int64                  `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID    int64                  `json:"category_id" binding:"required" gorm:"not null"`
	Name          string                 `json:"name" binding:"required" gorm:"type:varchar(200);not null"`
	Slug          string                 `json:"slug" gorm:"type:varchar(200);not null"`
	Description   string                 `json:"description" binding:"required" gorm:"type:text;not null"`
	BasePrice     float64                `json:"base_price" gorm:"type:numeric(12,2)"`
	IsActive      bool                   `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time              `json:"created_at" gorm:"autoCreateTimeUTC"`
	UpdatedAt     time.Time              `json:"updated_at" gorm:"autoUpdateTimeUTC"`
	Media         []ServiceMedia         `json:"media" gorm:"foreignKey:ServiceID"` // slice untuk media
	Spesification []ServiceSpesification `json:"spesification" gorm:"foreignKey:ServiceID"`
}

type ServiceMedia struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceID int64     `json:"service_id" gorm:"not null;index"`
	URL       string    `json:"url" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"type:varchar(50);not null"` // icon, thumbnail, gallery
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTimeUTC"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTimeUTC"`
}

type ServiceSpesification struct {
	ID                 int64                       `json:"id" gorm:"primaryKey"`
	ServiceID          int64                       `json:"service_id"  gorm:"not null;index"`
	Name               string                      `json:"name"`
	InputType          string                      `json:"input_type"`
	Options            datatypes.JSON              `json:"options"`
	IsRequired         bool                        `json:"is_required" gorm:"default:true"`
	IsActive           bool                        `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time                   `json:"created_at" gorm:"autoCreateTimeUTC"`
	UpdatedAt          time.Time                   `json:"updated_at" gorm:"autoUpdateTimeUTC"`
	SpesificationValue []ServiceSpesificationValue `json:"spesification_value" gorm:"foreignKey:SpesificationID;references:ID"`
}

type ServiceSpesificationValue struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	ServiceID       int64     `json:"service_id"  gorm:"not null;index"`
	SpesificationID int64     `json:"spesification_id"  gorm:"not null;index"`
	Value           string    `json:"value"`
	AdditionalPrice float64   `json:"additional_price" gorm:"type:numeric(12,2)"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTimeUTC"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTimeUTC"`
}

type RequestService struct {
	CategoryID  int64   `form:"category_id" binding:"required,gt=0"`
	Name        string  `form:"name" binding:"required,max=200"`
	Description string  `form:"description" binding:"required"`
	BasePrice   float64 `form:"base_price" binding:"gte=0"`
}

type RequestServiceSpesification struct {
	ID         int64          `json:"id"`
	ServiceId  int64          `json:"service_id" binding:"required,gt=0"`
	Name       string         `json:"name" binding:"required,min=3,max=100"`
	InputType  string         `json:"input_type" binding:"required,oneof=select boolean text number"`
	Options    datatypes.JSON `json:"options"`
	IsRequired *bool          `json:"is_required" binding:"required"`
}

type RequestServiceSpesificationValue struct {
	ServiceID       int64   `json:"service_id" binding:"required,gt=0"`
	SpesificationID int64   `json:"spesification_id" binding:"required,gt=0"`
	Value           string  `json:"value"`
	AdditionalPrice float64 `json:"additional_price" binding:"gte=0"`
}

type RequestAddServiceMedia struct {
	Type string `form:"type" binding:"required,oneof=gallery icon thumbnail"`
	URL  string `json:"-"`
}

type RequestUpdateServiceSpesificationValue struct {
	Value           string  `json:"value"`
	AdditionalPrice float64 `json:"additional_price" binding:"gte=0"`
}
