package repository

import (
	"context"
	"errors"
	"jasamc/models"
	"time"

	"gorm.io/gorm"
)

func (r *JasaRepository) GetAllServices(ctx context.Context) ([]models.Service, error) {
	var services []models.Service

	err := r.Database.WithContext(ctx).
		Preload("Media").                            // preload semua media
		Preload("Spesification").                    // preload spesifikasi
		Preload("Spesification.SpesificationValue"). // preload nilai spesifikasi
		Find(&services).Error
	if err != nil {
		return nil, err
	}

	// pastikan slice media dan spesifikasi tidak nil
	for i := range services {
		if services[i].Media == nil {
			services[i].Media = []models.ServiceMedia{}
		}
		if services[i].Spesification == nil {
			services[i].Spesification = []models.ServiceSpesification{}
		} else {
			// pastikan nested SpesificationValue tidak nil
			for j := range services[i].Spesification {
				if services[i].Spesification[j].SpesificationValue == nil {
					services[i].Spesification[j].SpesificationValue = []models.ServiceSpesificationValue{}
				}
			}
		}
	}

	return services, nil
}

// get category jasa by name
func (r *JasaRepository) FindServiceByName(ctx context.Context, name string) (*models.Service, error) {
	var jasa models.Service

	err := r.Database.WithContext(ctx).
		Table("services").
		Where("name = ?", name).
		First(&jasa).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &jasa, nil
}

// Get by slug
func (r *JasaRepository) FindServiceBySlug(ctx context.Context, slug string) (*models.Service, error) {
	var jasa models.Service
	err := r.Database.WithContext(ctx).
		Table("services").
		Where("slug = ?", slug).
		First(&jasa).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &jasa, nil
}

// Get jasa by ID
func (r *JasaRepository) FindServiceByID(ctx context.Context, id int64) (*models.Service, error) {
	var jasa models.Service
	err := r.Database.WithContext(ctx).
		Preload("Media").
		Preload("Spesification").
		Preload("Spesification.SpesificationValue").
		Table("services").
		Where("id = ?", id).
		First(&jasa).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &jasa, nil
}

// get service_media by id
func (r *JasaRepository) FindMediaServiceByID(ctx context.Context, id int64) (*models.ServiceMedia, error) {
	var media models.ServiceMedia
	err := r.Database.WithContext(ctx).
		Table("service_media").
		Where("id = ?", id).
		First(&media).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &media, nil
}

func (r *JasaRepository) CreateServiceWithMedia(
	ctx context.Context,
	service *models.Service,
) error {

	tx := r.Database.WithContext(ctx).Begin()

	if err := tx.Create(service).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *JasaRepository) DeleteMediaBYServiceIDandMediaID(
	ctx context.Context,
	serviceID, mediaID int64,
) error {

	result := r.Database.
		WithContext(ctx).
		Where("id = ? AND service_id = ?", mediaID, serviceID).
		Delete(&models.ServiceMedia{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// update status service
func (r *JasaRepository) ToggleServiceActiveAndReturn(
	ctx context.Context,
	id int64,
) (*models.Service, error) {

	// 1. Toggle is_active (atomic)
	result := r.Database.
		WithContext(ctx).
		Exec(`
			UPDATE services
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
	var service models.Service
	err := r.Database.
		WithContext(ctx).
		Table("services").
		Where("id = ?", id).
		First(&service).
		Error

	if err != nil {
		return nil, err
	}

	return &service, nil
}

// create spesification service
func (r *JasaRepository) InsertServiceSpesification(ctx context.Context, param *models.RequestServiceSpesification) (int64, error) {
	err := r.Database.WithContext(ctx).Table("service_spesifications").Create(param).Error

	if err != nil {
		return 0, err
	}

	return param.ID, nil
}

// find service spesification
func (r *JasaRepository) FindServiceSpesificationByName(
	ctx context.Context,
	serviceID int64,
	name string,
) (*models.ServiceSpesification, error) {

	var spec models.ServiceSpesification

	err := r.Database.WithContext(ctx).
		Table("service_spesifications").
		Where("service_id = ? AND name ILIKE ?", serviceID, name).
		First(&spec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &spec, nil
}

func (r *JasaRepository) FindServiceSpesificationByID(
	ctx context.Context,
	id int64,
) (*models.ServiceSpesification, error) {

	var ServiceSpec models.ServiceSpesification
	err := r.Database.WithContext(ctx).
		Table("service_spesifications").
		Where("id = ?", id).
		First(&ServiceSpec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ServiceSpec, nil
}

func (r *JasaRepository) AddNewServiceMedia(
	ctx context.Context,
	serviceID int64,
	param models.RequestAddServiceMedia,
) error {

	media := models.ServiceMedia{
		ServiceID: serviceID,
		Type:      param.Type,
		URL:       param.URL,
	}

	err := r.Database.WithContext(ctx).
		Table("service_media").
		Create(&media).Error

	if err != nil {
		return err
	}

	return nil
}

// cek data type di service ada berapa
func (r *JasaRepository) CountServiceMediaByType(
	ctx context.Context,
	serviceID int64,
	mediaType string,
) (int64, error) {

	var count int64

	err := r.Database.WithContext(ctx).
		Table("service_media").
		Where("service_id = ? AND type = ?", serviceID, mediaType).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *JasaRepository) InsertNewServiceSpesificationValue(
	ctx context.Context,
	param models.RequestServiceSpesificationValue,
) error {

	err := r.Database.WithContext(ctx).
		Table("service_spesification_values").
		Create(&param).Error

	if err != nil {
		return err
	}

	return nil
}

// cari value spesification dengan name
func (r *JasaRepository) FindServiceSpesificationValueByServiceAndSpec(
	ctx context.Context,
	serviceID int64,
	spesificationID int64,
	value string,
) (*models.ServiceSpesificationValue, error) {

	var specValue models.ServiceSpesificationValue

	err := r.Database.WithContext(ctx).
		Table("service_spesification_values").
		Where(
			"service_id = ? AND spesification_id = ? AND value ILIKE ?",
			serviceID,
			spesificationID,
			value,
		).
		First(&specValue).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &specValue, nil
}

func (r *JasaRepository) FindServiceSpesificationValueByServiceIDBySpecIDByValueID(
	ctx context.Context,
	serviceID int64,
	specID int64,
	valueID int64,
) (*models.ServiceSpesificationValue, error) {

	var value models.ServiceSpesificationValue

	err := r.Database.WithContext(ctx).
		Table("service_spesification_values").
		Where(
			"service_id = ? AND spesification_id = ? AND id = ?",
			serviceID,
			specID,
			valueID,
		).
		First(&value).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &value, nil
}

func (r *JasaRepository) UpdateServiceSpesificationValue(
	ctx context.Context,
	serviceID int64,
	specID int64,
	valueID int64,
	param models.RequestUpdateServiceSpesificationValue,
) error {

	result := r.Database.WithContext(ctx).
		Table("service_spesification_values").
		Where(
			"service_id = ? AND spesification_id = ? AND id = ?",
			serviceID, specID, valueID,
		).
		Updates(map[string]interface{}{
			"value":            param.Value,
			"additional_price": param.AdditionalPrice,
			"updated_at":       time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("spesification value tidak ditemukan")
	}

	return nil
}

func (r *JasaRepository) DeleteServiceSpesification(
	ctx context.Context,
	serviceID int64,
	specID int64,
) error {

	result := r.Database.WithContext(ctx).
		Where("service_id = ? AND id = ?", serviceID, specID).
		Delete(&models.ServiceSpesification{})

	if result.RowsAffected == 0 {
		return errors.New("spesification tidak ditemukan")
	}

	return result.Error
}

// delete services
func (r *JasaRepository) DeleteServiceByID(ctx context.Context, serviceID int64) error {
	result := r.Database.WithContext(ctx).
		Delete(&models.Service{}, serviceID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// toggel active spec
// update status service spesification
func (r *JasaRepository) ToggleServiceSpesificationActiveAndReturn(
	ctx context.Context,
	serviceID int64,
	spesificationID int64,
) (*models.ServiceSpesification, error) {

	// 1. Toggle is_active (atomic)
	result := r.Database.
		WithContext(ctx).
		Exec(`
			UPDATE service_spesifications
			SET is_active = NOT is_active
			WHERE service_id = ? AND id = ?
		`, serviceID, spesificationID)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// 2. Ambil data terbaru
	var spesification models.ServiceSpesification
	err := r.Database.
		WithContext(ctx).
		Table("service_spesifications").
		Where("service_id = ? AND id = ?", serviceID, spesificationID).
		First(&spesification).
		Error

	if err != nil {
		return nil, err
	}

	return &spesification, nil
}

// update status service spesification
func (r *JasaRepository) ToggleServiceSpesificationRequiredAndReturn(
	ctx context.Context,
	serviceID int64,
	spesificationID int64,
) (*models.ServiceSpesification, error) {

	// 1. Ambil status saat ini
	var current models.ServiceSpesification
	err := r.Database.
		WithContext(ctx).
		Select("id, is_active, is_required").
		Where("service_id = ? AND id = ?", serviceID, spesificationID).
		First(&current).
		Error

	if err != nil {
		return nil, err
	}

	// 2. VALIDASI:
	// jika saat ini inactive & required=false
	// maka toggle akan jadi required=true → TOLAK
	if !current.IsActive && !current.IsRequired {
		return nil, errors.New(
			"spesification tidak boleh required jika tidak aktif",
		)
	}

	// 3. Toggle is_required
	result := r.Database.
		WithContext(ctx).
		Exec(`
			UPDATE service_spesifications
			SET is_required = NOT is_required
			WHERE service_id = ? AND id = ?
		`, serviceID, spesificationID)

	if result.Error != nil {
		return nil, result.Error
	}

	var spesification models.ServiceSpesification
	err = r.Database.
		WithContext(ctx).
		Where("service_id = ? AND id = ?", serviceID, spesificationID).
		First(&spesification).
		Error

	return &spesification, err
}
