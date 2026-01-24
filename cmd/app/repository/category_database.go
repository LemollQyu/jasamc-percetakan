package repository

import (
	"context"
	"errors"
	"jasamc/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// query category db

// get by id
func (r *JasaRepository) FindCategoryJasaByID(ctx context.Context, jasaID int64) (*models.CategoryJasa, error) {
	var categoryJasa models.CategoryJasa
	err := r.Database.WithContext(ctx).Table("service_categories").Where("id = ?", jasaID).First(&categoryJasa).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &categoryJasa, nil

}

// get category jasa by name
func (r *JasaRepository) FindCategoryJasaByName(ctx context.Context, name string) (*models.CategoryJasa, error) {
	var categoryJasa models.CategoryJasa

	err := r.Database.WithContext(ctx).
		Table("service_categories").
		Where("name = ?", name).
		First(&categoryJasa).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &categoryJasa, nil
}

// Get by slug
func (r *JasaRepository) FindCategoryJasaBySlug(ctx context.Context, slug string) (*models.CategoryJasa, error) {
	var categoryJasa models.CategoryJasa
	err := r.Database.WithContext(ctx).
		Table("service_categories").
		Where("slug = ?", slug).
		First(&categoryJasa).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &categoryJasa, nil
}

// insert category jasa
func (r *JasaRepository) InsertCategoryJasa(ctx context.Context, categoryJasa *models.ParamCreateCategoryJasa) (int64, error) {
	err := r.Database.WithContext(ctx).Table("service_categories").Create(categoryJasa).Error

	if err != nil {
		return 0, err
	}

	return categoryJasa.ID, nil
}

// delete category
func (r *JasaRepository) DeleteCategoryJasa(ctx context.Context, id int64) error {
	err := r.Database.WithContext(ctx).
		Table("service_categories").
		Where("id = ?", id).
		Delete(nil).Error

	if err != nil {
		return err
	}

	return nil
}

// cek category jasa ada fk jasanya tidak
func (r *JasaRepository) CountServicesByCategoryID(ctx context.Context, categoryID int64) (int64, error) {
	var count int64
	err := r.Database.WithContext(ctx).
		Table("services").
		Where("category_id = ?", categoryID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// get all category jasa
func (r *JasaRepository) GetAllCategoryJasa(ctx context.Context) ([]models.CategoryJasa, error) {
	var categories []models.CategoryJasa

	err := r.Database.
		WithContext(ctx).
		Table("service_categories").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// update meta dengan icon
func (r *JasaRepository) UpdateCategoryMeta(
	ctx context.Context,
	categoryID int64,
	meta []byte,
) error {
	return r.Database.
		WithContext(ctx).
		Table("service_categories").
		Where("id = ?", categoryID).
		Update("meta", datatypes.JSON(meta)).
		Error
}

// update categirynya (ini masih bingung di implemetasikan atau tidak)

// update setstatus caetgory
func (r *JasaRepository) ToggleCategoryJasaActiveAndReturn(
	ctx context.Context,
	id int64,
) (*models.CategoryJasa, error) {

	// 1. Toggle is_active (atomic)
	result := r.Database.
		WithContext(ctx).
		Exec(`
			UPDATE service_categories
			SET is_active = NOT is_active
			WHERE id = ?
		`, id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// 2. Ambil data terbaru
	var category models.CategoryJasa
	err := r.Database.
		WithContext(ctx).
		Table("service_categories").
		Where("id = ?", id).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}
