package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type JasaRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewJasaRepository(db *gorm.DB, redis *redis.Client) *JasaRepository {
	return &JasaRepository{
		Database: db,
		Redis:    redis,
	}
}
