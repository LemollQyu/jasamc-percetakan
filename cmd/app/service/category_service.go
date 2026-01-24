package service

import (
	"context"
	"jasamc/infrastructure/log"
	"jasamc/models"
	"time"

	"github.com/sirupsen/logrus"
)

// cari category by id
// untuk database
func (s *JasaService) GetCategoryJasaByIDFromDB(ctx context.Context, id int64) (*models.CategoryJasa, error) {
	jasa, err := s.JasaRepo.FindCategoryJasaByID(ctx, id)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"jasa_id": id,
			"error":   err.Error(),
		}).Error("gagal di service s.JasaRepo.FindJasaByID")
		return nil, err
	}
	return jasa, nil
}

// unutk cache
func (s *JasaService) GetCategoryJasaByIDFromRead(
	ctx context.Context,
	id int64,
) (*models.CategoryJasa, error) {

	// Redis GET
	category, err := s.JasaRepo.GetCategoryJasaByIDFromRedis(ctx, id)
	if err != nil {
		log.Logger.Error("REDIS GET ERROR (by id)", err)
	}

	if category != nil {
		log.Logger.WithField("id", id).Info("CACHE HIT category_jasa:id")
		return category, nil
	}

	// DB GET
	category, err = s.JasaRepo.FindCategoryJasaByID(ctx, id)

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"category_id": id,
			"error":       err.Error(),
		}).Error("DB ERROR FindCategoryJasaByID")
		return nil, err
	}

	if category == nil {
		return nil, nil
	}

	log.Logger.WithField("id", id).Info("CACHE MISS category_jasa:id")

	// Redis SET (ASYNC)
	go func(data *models.CategoryJasa) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.JasaRepo.SetCategoryJasaByIDToRedis(bgCtx, data); err != nil {
			log.Logger.Error("REDIS SET ERROR (by id)", err)
		}
	}(category)

	return category, nil
}

func (s *JasaService) CreateNewCategoryJasa(ctx context.Context, param *models.ParamCreateCategoryJasa) (int64, error) {
	categoryJasaID, err := s.JasaRepo.InsertCategoryJasa(ctx, param)

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
			"error": err.Error(),
		}).Error("gagal di service s.JasaRepo.InsertCategoryJasa")
		return 0, err
	}

	// 2. Hapus cache Redis (agar GET berikutnya ambil dari DB)
	if err := s.JasaRepo.DeleteAllCategoryCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis category_jasa:all", err)
	}

	return categoryJasaID, nil

}

func (s *JasaService) FindCategoryJasaByName(ctx context.Context, name string) (*models.CategoryJasa, error) {
	categoryJasa, err := s.JasaRepo.FindCategoryJasaByName(ctx, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":  name,
			"error": err.Error(),
		}).Error("gagal di service s.JasaRepo.FindCategoryJasaByName")
		return nil, err
	}

	return categoryJasa, nil
}

func (s *JasaService) FindCategoryJasaBySlug(ctx context.Context, slug string) (*models.CategoryJasa, error) {
	categoryJasa, err := s.JasaRepo.FindCategoryJasaBySlug(ctx, slug)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":  slug,
			"error": err.Error(),
		}).Error("gagal di service s.JasaRepo.FindCategoryJasaBySlug")
		return nil, err
	}

	return categoryJasa, nil
}

func (s *JasaService) CountServiceByCategory(ctx context.Context, categoryID int64) (int64, error) {
	count, err := s.JasaRepo.CountServicesByCategoryID(ctx, categoryID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"category_id": categoryID,
			"error":       err.Error(),
		}).Error("gagal di service CountServicesByCategoryID")
		return 0, err
	}

	return count, nil
}

func (s *JasaService) DeleteCategoryJasa(ctx context.Context, id int64) error {
	err := s.JasaRepo.DeleteCategoryJasa(ctx, id)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"category_id": id,
			"error":       err.Error(),
		}).Error("gagal di service DeleteCategoryJasa")
		return err
	}

	// Hapus cache Redis (agar GET berikutnya ambil dari DB)
	if err := s.JasaRepo.DeleteCategoryJasaCacheByID(ctx, id); err != nil {
		log.Logger.Error("Gagal delete cache category_jasa:id", err)
	}

	if err := s.JasaRepo.DeleteAllCategoryCache(ctx); err != nil {
		log.Logger.Error("Gagal delete cache category_jasa:all", err)
	}

	return nil
}

func (s *JasaService) GetAllCategoryJasa(
	ctx context.Context,
) ([]models.CategoryJasa, error) {

	// 1. Redis GET
	categories, err := s.JasaRepo.GetAllCategoryJasaFromRedis(ctx)
	if err != nil {
		log.Logger.Error("REDIS GET ERROR", err)
	}

	if len(categories) > 0 {
		log.Logger.Info("CACHE HIT category_jasa:all")
		return categories, nil
	}

	// 2. DB
	categories, err = s.JasaRepo.GetAllCategoryJasa(ctx)
	if err != nil {
		return nil, err
	}

	log.Logger.Infof("CACHE MISS → DB RETURNED %d ROWS", len(categories))

	// 3. Redis SET ASYNC (DETACHED)
	go func(data []models.CategoryJasa) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.JasaRepo.SetAllCategoryJasaToRedis(bgCtx, data); err != nil {
			log.Logger.Error("REDIS SET ERROR (BG)", err)
		} else {
			log.Logger.Info("REDIS SET category_jasa:all SUCCESS (BG)")
		}
	}(categories)

	return categories, nil
}

func (s *JasaService) UpdateCategoryMeta(
	ctx context.Context,
	categoryID int64,
	meta []byte,
) error {
	err := s.JasaRepo.UpdateCategoryMeta(ctx, categoryID, meta)
	if err != nil {
		return err
	}
	// Hapus cache Redis (agar GET berikutnya ambil dari DB)
	// 2. DELETE CACHE (SYNC)

	if err := s.JasaRepo.DeleteCategoryJasaCacheByID(ctx, categoryID); err != nil {
		log.Logger.Error("Gagal delete cache category_jasa:id", err)
	}

	if err := s.JasaRepo.DeleteAllCategoryCache(ctx); err != nil {
		log.Logger.Error("Gagal delete cache category_jasa:all", err)
	}

	return nil
}

func (s *JasaService) UpdateStatusCategoryJasa(
	ctx context.Context,
	categoryID int64,
	isActive bool,
) (*models.CategoryJasa, error) {

	category, err := s.JasaRepo.ToggleCategoryJasaActiveAndReturn(ctx, categoryID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"category_id": categoryID,
			"error":       err.Error(),
		}).Error("gagal di service ToggleCategoryJasaActive")

		return nil, err
	}

	// INVALIDATE CACHE (WAJIB SYNC)
	if err := s.JasaRepo.DeleteCategoryJasaCacheByID(ctx, categoryID); err != nil {
		log.Logger.Error("Gagal menghapus cache category_jasa:id", err)
	}

	if err := s.JasaRepo.DeleteAllCategoryCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache category_jasa:all", err)
	}

	return category, nil
}
