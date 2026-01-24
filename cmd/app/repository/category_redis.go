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
	cacheKeyAllCategoryJasa  = "category_jasa:all"
	cacheKeyCategoryJasaInfo = "category_jasa:%d"
)

func (r *JasaRepository) GetAllCategoryJasaFromRedis(
	ctx context.Context,
) ([]models.CategoryJasa, error) {

	var categories []models.CategoryJasa

	categoryStr, err := r.Redis.Get(ctx, cacheKeyAllCategoryJasa).Result()
	if err != nil {
		if err == redis.Nil {
			// cache miss
			return []models.CategoryJasa{}, nil
		}

		return nil, err
	}

	// unmarshal json → struct
	err = json.Unmarshal([]byte(categoryStr), &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// SetAllCategoryJasaToRedis set all category jasa to redis.
func (r *JasaRepository) SetAllCategoryJasaToRedis(
	ctx context.Context,
	categories []models.CategoryJasa,
) error {

	categoryJSON, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	err = r.Redis.SetEx(
		ctx,
		cacheKeyAllCategoryJasa,
		categoryJSON,
		20*time.Minute,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *JasaRepository) DeleteAllCategoryCache(ctx context.Context) error {

	return r.Redis.Del(ctx, cacheKeyAllCategoryJasa).Err()
}

// helper key get category jasa id
func cacheKeyCategoryJasaByID(id int64) string {
	return fmt.Sprintf(cacheKeyCategoryJasaInfo, id)
}

// get category id di redis
func (r *JasaRepository) GetCategoryJasaByIDFromRedis(
	ctx context.Context,
	id int64,
) (*models.CategoryJasa, error) {

	var category models.CategoryJasa
	key := cacheKeyCategoryJasaByID(id)

	categoryStr, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // cache miss
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(categoryStr), &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// set category by id in redis
func (r *JasaRepository) SetCategoryJasaByIDToRedis(
	ctx context.Context,
	category *models.CategoryJasa,
) error {

	key := cacheKeyCategoryJasaByID(category.ID)

	data, err := json.Marshal(category)
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

// delete cetagory in redis
func (r *JasaRepository) DeleteCategoryJasaCacheByID(
	ctx context.Context,
	id int64,
) error {

	key := cacheKeyCategoryJasaByID(id)
	return r.Redis.Del(ctx, key).Err()
}
