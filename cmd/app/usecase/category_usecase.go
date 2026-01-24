package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"jasamc/infrastructure/log"
	"jasamc/models"
	"jasamc/utils"
	"mime/multipart"

	"github.com/sirupsen/logrus"
)

func (u *JasaUsecase) GetCategoryJasaByID(ctx context.Context, id int64) (*models.CategoryJasa, error) {
	categoryJasa, err := u.JasaService.GetCategoryJasaByIDFromRead(ctx, id)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"jasa_id": id,
			"error":   err.Error(),
		}).Error("JasaUsecase: Gagal get category jasa by id, u.JasaService.GetCategoryJasaByID")
		return nil, err
	}

	return categoryJasa, nil
}

func (u *JasaUsecase) GetCategoryJasaByName(ctx context.Context, name string) (*models.CategoryJasa, error) {
	categoryJasa, err := u.JasaService.FindCategoryJasaByName(ctx, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":  name,
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di s.JasaRepo.FindCategoryJasaByName")
		return nil, err
	}

	return categoryJasa, nil
}

func (u *JasaUsecase) GetCategoryJasaBySlug(ctx context.Context, slug string) (*models.CategoryJasa, error) {
	categoryJasa, err := u.JasaService.FindCategoryJasaBySlug(ctx, slug)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"slug":  slug,
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di s.JasaRepo.FindCategoryJasaBySlug")
		return nil, err
	}

	return categoryJasa, nil
}

func (u *JasaUsecase) CreateCategoryJasa(ctx context.Context, param *models.ParamCreateCategoryJasa) (int64, error) {

	// Normalisasi: meta kosong → nil
	if param.Meta != nil && len(param.Meta) == 0 {
		param.Meta = nil
	}

	// CEK NAME EXIST
	existName, err := u.JasaService.FindCategoryJasaByName(ctx, param.Name)
	if err != nil { // error sistem
		log.Logger.WithFields(logrus.Fields{
			"name":  param.Name,
			"error": err.Error(),
		}).Error("Gagal create category jasa by name")
		return 0, err
	}

	if existName != nil { // data ditemukan
		return 0, utils.ErrNameExists
	}

	// Generate slug
	param.Slug = utils.GenerateSlug(param.Name)

	// CEK SLUG EXIST
	existSlug, err := u.JasaService.FindCategoryJasaBySlug(ctx, param.Slug)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"slug":  param.Slug,
			"error": err.Error(),
		}).Error("Gagal create category jasa by slug")
		return 0, err
	}

	if existSlug != nil {
		return 0, utils.ErrSlugExists
	}

	// INSERT
	categoryJasaID, err := u.JasaService.CreateNewCategoryJasa(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
			"error": err.Error(),
		}).Error("Gagal insert category jasa")
		return 0, err
	}

	return categoryJasaID, nil
}

// usecase delete category jasa
func (u *JasaUsecase) DeleteCategoryJasa(ctx context.Context, id int64) error {

	// CEK CATEGORY EXIST?
	category, err := u.JasaService.GetCategoryJasaByIDFromDB(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category jasa tidak ditemukan")
	}

	// CEK APAKAH CATEGORY DIPAKAI OLEH TABLE services
	count, err := u.JasaService.CountServiceByCategory(ctx, id)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("category tidak bisa dihapus karena masih digunakan oleh service lain")
	}

	// CALL SERVICE UNTUK DELETE
	err = u.JasaService.DeleteCategoryJasa(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// all category jasa
func (u *JasaUsecase) GetAllCategoryJasa(ctx context.Context) ([]models.CategoryJasa, error) {
	categories, err := u.JasaService.GetAllCategoryJasa(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// update meta dengan icon
func (u *JasaUsecase) UpdateCategoryIcon(
	ctx context.Context,
	categoryID int64,
	file *multipart.FileHeader,
) error {

	// cek category ada nggak
	category, err := u.JasaService.GetCategoryJasaByIDFromDB(ctx, categoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category jasa tidak ditemukan")
	}

	// 1. upload icon
	iconURL, err := u.StorageService.UploadCategoryIcon(ctx, file)
	if err != nil {
		return err
	}

	// 2. build meta
	meta := map[string]interface{}{
		"icon": iconURL,
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	// 3. update DB
	return u.JasaService.UpdateCategoryMeta(ctx, categoryID, metaJSON)
}

// usecase set status category jasa
func (u *JasaUsecase) UpdateStatusCategoryJasa(ctx context.Context, id int64) (*models.CategoryJasa, error) {
	category, err := u.JasaService.GetCategoryJasaByIDFromDB(ctx, id)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category jasa tidak ditemukan")
	}

	category, err = u.JasaService.UpdateStatusCategoryJasa(ctx, id, category.IsActive)
	if err != nil {
		return nil, err
	}

	return category, nil

}
