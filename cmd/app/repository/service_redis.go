package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"jasamc/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	cacheKeyAllService  = "services:all"
	cacheKeyServiceInfo = "services:%d"
)

//
// ==========================
// GET ALL SERVICES
// ==========================
//

// GetAllServicesFromRedis get all services from redis
func (r *JasaRepository) GetAllServicesFromRedis(
	ctx context.Context,
) ([]models.Service, error) {

	var services []models.Service

	serviceStr, err := r.Redis.Get(ctx, cacheKeyAllService).Result()
	if err != nil {
		if err == redis.Nil {
			// cache miss
			return []models.Service{}, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(serviceStr), &services); err != nil {
		return nil, err
	}

	return services, nil
}

// SetAllServicesToRedis set all services to redis
func (r *JasaRepository) SetAllServicesToRedis(
	ctx context.Context,
	services []models.Service,
) error {

	data, err := json.Marshal(services)
	if err != nil {
		return err
	}

	return r.Redis.SetEx(
		ctx,
		cacheKeyAllService,
		data,
		15*time.Minute,
	).Err()
}

// DeleteAllServiceCache delete cache all services
func (r *JasaRepository) DeleteAllServiceCache(ctx context.Context) error {
	return r.Redis.Del(ctx, cacheKeyAllService).Err()
}

//
// ==========================
// GET SERVICE BY ID
// ==========================
//

// helper redis key
func cacheKeyServiceByID(id int64) string {
	return fmt.Sprintf(cacheKeyServiceInfo, id)
}

// GetServiceByIDFromRedis get service by id from redis
func (r *JasaRepository) GetServiceByIDFromRedis(
	ctx context.Context,
	id int64,
) (*models.Service, error) {

	var service models.Service
	key := cacheKeyServiceByID(id)

	serviceStr, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // cache miss
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(serviceStr), &service); err != nil {
		return nil, err
	}

	return &service, nil
}

// SetServiceByIDToRedis set service by id to redis
func (r *JasaRepository) SetServiceByIDToRedis(
	ctx context.Context,
	service *models.Service,
) error {

	key := cacheKeyServiceByID(service.ID)

	data, err := json.Marshal(service)
	if err != nil {
		return err
	}

	return r.Redis.SetEx(
		ctx,
		key,
		data,
		5*time.Minute,
	).Err()
}

// DeleteServiceCacheByID delete service cache by id
func (r *JasaRepository) DeleteServiceCacheByID(
	ctx context.Context,
	id int64,
) error {

	key := cacheKeyServiceByID(id)
	return r.Redis.Del(ctx, key).Err()
}
