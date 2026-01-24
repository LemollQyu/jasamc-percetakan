package usecase

import (
	"context"
	"errors"
	"fmt"
	"jasamc/infrastructure/log"
	"jasamc/models"
	"jasamc/utils"
	"mime/multipart"

	"github.com/sirupsen/logrus"
)

func (u *JasaUsecase) GetAllJasa(ctx context.Context) ([]models.Service, error) {
	// Panggil service
	categories, err := u.JasaService.GetAllServices(ctx) // pastikan service return []models.CategoryJasa
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get all kategori jasa in usecase")
		return nil, err
	}

	return categories, nil
}

func (u *JasaUsecase) GetServiceByName(ctx context.Context, name string) (*models.Service, error) {
	jasa, err := u.JasaService.GetServiceByName(ctx, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.GetServiceByName")
		return nil, err
	}

	return jasa, nil
}

func (u *JasaUsecase) GetServiceBySlug(ctx context.Context, name string) (*models.Service, error) {
	jasa, err := u.JasaService.GetServiceBySlug(ctx, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.GetServiceBySlug")
		return nil, err
	}

	return jasa, nil
}

func (u *JasaUsecase) CreateAndUploadFileService(
	ctx context.Context,
	param models.RequestService,
	files map[string][]*multipart.FileHeader,
) (*models.Service, error) {

	// validasi category_id memang ada atau nggak
	category, err := u.JasaService.GetCategoryJasaByIDFromDB(ctx, param.CategoryID)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category tidak ditemukan")
	}

	// validasi name dan slug
	jasa, err := u.JasaService.GetServiceByName(ctx, param.Name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.GetServiceByName")
		return nil, err
	}

	if jasa != nil {
		return nil, errors.New("name sudah dipakai")
	}

	jasa, err = u.JasaService.GetServiceBySlug(ctx, utils.GenerateSlug(param.Name))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.GetServiceByName")
		return nil, err
	}

	if jasa != nil {
		return nil, errors.New("slug sudah dipakai")
	}

	// VALIDASI FILE INPUT

	// Thumbnail WAJIB
	thumbnails, hasThumbnail := files["thumbnail"]
	if !hasThumbnail || len(thumbnails) == 0 {
		return nil, errors.New("thumbnail wajib diisi")
	}

	// Maksimal 4 thumbnail
	if len(thumbnails) > 4 {
		return nil, errors.New("maksimal thumbnail 4")
	}

	// Gallery WAJIB
	galleries, hasGallery := files["gallery"]
	if !hasGallery || len(galleries) == 0 {
		return nil, errors.New("gallery wajib diisi")
	}

	// Icon opsional
	icons := files["icon"]

	// PREPARE URL HOLDER

	var (
		iconURL       string
		thumbnailURLs []string
		galleryURLs   []string
	)

	// UPLOAD ICON (OPSIONAL)

	if len(icons) > 0 {
		url, err := u.StorageService.UploadServiceIcon(ctx, icons[0])
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("jasaUsecase: gagal di u.JasaService.UploadServiceIcon")
			return nil, err
		}
		iconURL = url
	}

	//  UPLOAD THUMBNAIL

	for _, file := range thumbnails {
		url, err := u.StorageService.UploadServiceThumbnail(ctx, file)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("jasaUsecase: gagal di u.JasaService.UploadServiceThumbnail")
			return nil, err
		}
		thumbnailURLs = append(thumbnailURLs, url)
	}

	// UPLOAD GALLERY

	for _, file := range galleries {
		url, err := u.StorageService.UploadServiceGallery(ctx, file)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("jasaUsecase: gagal di u.JasaService.UploadSeviceGallery")
			return nil, err
		}
		galleryURLs = append(galleryURLs, url)
	}

	//  BUILD MEDIA MODELS

	var medias []models.ServiceMedia

	if iconURL != "" {
		medias = append(medias, models.ServiceMedia{
			URL:  iconURL,
			Type: "icon",
		})
	}

	for _, url := range thumbnailURLs {
		medias = append(medias, models.ServiceMedia{
			URL:  url,
			Type: "thumbnail",
		})
	}

	for _, url := range galleryURLs {
		medias = append(medias, models.ServiceMedia{
			URL:  url,
			Type: "gallery",
		})
	}

	// BUILD SERVICE MODEL

	service := &models.Service{
		CategoryID:  param.CategoryID,
		Name:        param.Name,
		Slug:        utils.GenerateSlug(param.Name),
		Description: param.Description,
		BasePrice:   param.BasePrice,
		Media:       medias,
	}

	// =========================
	// 8️⃣ SAVE TO DB
	// =========================
	if err := u.JasaService.CreateService(ctx, service); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.CreateService")
		return nil, err
	}

	return service, nil
}

func (u *JasaUsecase) AddServiceMediaByType(
	ctx context.Context,
	serviceID int64,
	param models.RequestAddServiceMedia,
	file *multipart.FileHeader,
) error {

	// 1️⃣ validasi service
	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {
		return err
	}
	if service == nil {
		return errors.New("service tidak ditemukan")
	}

	// 2️⃣ validasi jumlah media berdasarkan type
	count, err := u.JasaService.CountServiceMediaByType(ctx, serviceID, param.Type)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": fmt.Sprintf("%d, %s", serviceID, param.Type),
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.CountServiceMediaByType")
		return err

	}

	switch param.Type {
	case "icon":
		if count >= 1 {
			return errors.New("icon sudah ada")
		}
	case "thumbnail":
		if count >= 3 {
			return errors.New("thumbnail maksimal 3")
		}
	case "gallery":
		// bebas
	default:
		return errors.New("file tidak sesuai dengan type")
	}

	var (
		url string
	)

	// 3️⃣ upload sesuai type
	switch param.Type {
	case "icon":
		url, err = u.StorageService.UploadServiceIcon(ctx, file)
	case "thumbnail":
		url, err = u.StorageService.UploadServiceThumbnail(ctx, file)
	case "gallery":
		url, err = u.StorageService.UploadServiceGallery(ctx, file)
	}

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.StorageService")
		return err
	}

	// 4️⃣ insert ke DB
	err = u.JasaService.AddServiceMedia(ctx, serviceID, models.RequestAddServiceMedia{
		Type: param.Type,
		URL:  url,
	})

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"type":  param.Type,
			"url":   url,
			"error": err.Error(),
		}).Error("jasaUsecase: gagal di u.JasaService.AddServiceMedia")
		return err
	}

	return nil
}

func (u *JasaUsecase) DeleteMediaInServiceFromDBAndStorage(
	ctx context.Context,
	serviceID, mediaID int64,
) error {

	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {
		log.Logger.WithError(err).
			Error("jasaUsecase: u.JasaService.GetServiceByIDFromDB")
		return err
	}
	if service == nil {
		return errors.New("service tidak ditemukan")
	}

	media, err := u.JasaService.GetServiceMediaByID(ctx, mediaID)
	if err != nil {
		log.Logger.WithError(err).
			Error("jasaUsecase: u.JasaService.GetServiceMediaByID")
		return err
	}
	if media == nil {
		return errors.New("media tidak ditemukan")
	}

	// VALIDASI RELASI
	if media.ServiceID != serviceID {
		return errors.New("media tidak terkait dengan service ini")
	}

	// DELETE STORAGE DULU
	if err := u.StorageService.DeleteMediaByURL(ctx, media.URL); err != nil {
		log.Logger.WithError(err).
			Error("jasaUsecase: u.JasaService.DeleteMediaByURL")
		return err
	}

	// BARU DELETE DB
	if err := u.JasaService.DeleteMediaBYServiceIDandMediaID(ctx, serviceID, mediaID); err != nil {
		log.Logger.WithError(err).
			Error("jasaUsecase: u.JasaService.DeleteMediaBYServiceIDandMediaID")
		return err
	}

	return nil
}

func (u *JasaUsecase) SetStatusService(ctx context.Context, serviceID int64) (*models.Service, error) {
	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {
		log.Logger.
			Error("jasaUsecase: u.JasaService.GetServiceByIDFromDB")

		return nil, err
	}

	if service == nil {
		return nil, errors.New("service tidak ditemukan")
	}

	service, err = u.JasaService.UpdateStatusService(ctx, serviceID)
	if err != nil {
		log.Logger.
			Error("jasaUsecase: u.JasaService.UpdateStatusService")
		return nil, err
	}

	return service, nil
}

func (u *JasaUsecase) CreateServiceSpesification(ctx context.Context, param *models.RequestServiceSpesification) (int64, error) {

	service, err := u.JasaService.GetServiceByIDFromDB(ctx, param.ServiceId)
	if err != nil {

		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceByIDFromDB")
		return 0, err
	}

	if service == nil {
		return 0, errors.New("service tidak ditemukan")
	}

	serviceSpec, err := u.JasaService.GetServiceSpesificationByName(ctx, param.ServiceId, param.Name)
	if err != nil {
		log.Logger.Error("jasaService: u.JasaService.GetServiceSpesificationByName")
		return 0, err
	}

	if serviceSpec != nil {
		return 0, errors.New("name sudah digunakan")
	}

	specID, err := u.JasaService.CreateServiceSpesification(ctx, param)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.CreateServiceSpesification")
		return 0, nil
	}

	return specID, nil
}

func (u *JasaUsecase) CreateServiceSpesificationValue(ctx context.Context, param models.RequestServiceSpesificationValue) error {

	// validasi service ada nggak
	service, err := u.JasaService.GetServiceByIDFromDB(ctx, param.ServiceID)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceByIDFromDB")
		return err
	}

	if service == nil {
		return errors.New("service tidak ditemukan")
	}

	// validasi spesification ada nggak
	serviceSpec, err := u.JasaService.GetServiceSpesificationByID(ctx, param.SpesificationID)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceSpesificationByID")
		return err
	}

	if serviceSpec == nil {
		return errors.New("service spesification tidak ditemukan")
	}

	exist, err := u.JasaService.GetServiceSpesificationValueByName(ctx, param.ServiceID, param.SpesificationID, param.Value)
	if err != nil {
		log.Logger.Error("jasaHandler: u.JasaService.GetServiceSpesificationValueByName")
		return err
	}

	if exist != nil {
		return errors.New("value spesification sudah digunakan / dibuat")
	}

	err = u.JasaService.CreateServiceSpesificationValue(ctx, param)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.CreateServiceSpesificationValue")
		return err
	}

	return nil
}

// ini nggak kepake tapi nanti entahlah, buat ambil specValue
func (u *JasaUsecase) GetServiceSpesificationValueByServiceIDBySpecIDByValueID(ctx context.Context, serviceID, specID, valueID int64) (*models.ServiceSpesificationValue, error) {
	specValue, err := u.JasaService.GetServiceSpesificationValueByServiceIDBySpecIDByValueID(ctx, serviceID, specID, valueID)

	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceSpesificationValueByServiceIDBySpecIDByValueID")
		return nil, err
	}

	return specValue, nil
}

func (u *JasaUsecase) UpdateServiceSpesificationValue(
	ctx context.Context,
	serviceID, specID, valueID int64,
	param models.RequestUpdateServiceSpesificationValue,
) error {
	specValue, err := u.JasaService.GetServiceSpesificationValueByServiceIDBySpecIDByValueID(ctx, serviceID, specID, valueID)

	if err != nil {
		log.Logger.Error("JasaUsecase: u.JasaService.GetServiceSpesificationValueByServiceIDBySpecIDByValueID")
		return err
	}

	if specValue == nil {
		return errors.New("service spesification value tidak ditemukan")
	}

	err = u.JasaService.UpdateServiceSpesificationValue(ctx, serviceID, specID, valueID, param)
	if err != nil {
		log.Logger.Error("JasaUsecase: u.JasaService.UpdateServiceSpesificationValue")
		return err
	}

	return nil

}

func (u *JasaUsecase) DeleteServiceSpesification(ctx context.Context, serviceID, specID int64) error {
	// cek id
	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {

		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceByIDFromDB")
		return err
	}

	if service == nil {
		return errors.New("service tidak ditemukan")
	}

	serviceSpec, err := u.JasaService.GetServiceSpesificationByID(ctx, specID)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceSpesificationByID")
		return err
	}

	if serviceSpec == nil {
		return errors.New("service spesification tidak ditemukan")
	}

	err = u.JasaService.DeleteServiceSpesification(ctx, serviceID, specID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"specID":    specID,
		})
		return err
	}

	return nil

}

func (u *JasaUsecase) DeleteServiceByID(ctx context.Context, serviceID int64) error {
	// 1️⃣ Ambil service + media
	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {
		log.Logger.Error("jasaUsecase: GetServiceByIDFromDB")
		return err
	}

	if service == nil {
		return errors.New("service tidak ditemukan")
	}

	// 2️⃣ Hapus semua file di storage
	for _, media := range service.Media {
		err := u.StorageService.DeleteMediaByURL(ctx, media.URL)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"serviceID": serviceID,
				"mediaURL":  media.URL,
				"error":     err.Error(),
			}).Error("gagal hapus media storage")

			// STOP → jangan hapus DB
			return err
		}
	}

	// 3️⃣ Hapus service di DB (CASCADE)
	err = u.JasaService.DeleteServiceByID(ctx, serviceID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"error":     err.Error(),
		}).Error("gagal hapus service di DB")

		return err
	}

	return nil
}

func (u *JasaUsecase) ToggleServiceSpesificationStatus(ctx context.Context, serviceID, specID int64) (*models.ServiceSpesification, error) {

	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {
		log.Logger.Error("jasaUsecase: GetServiceByIDFromDB")
		return nil, err
	}

	if service == nil {
		return nil, errors.New("service tidak ditemukan")
	}

	serviceSpec, err := u.JasaService.GetServiceSpesificationByID(ctx, specID)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceSpesificationByID")
		return nil, err
	}

	if serviceSpec == nil {
		return nil, errors.New("service spesification tidak ditemukan")
	}

	serviceSpec, err = u.JasaService.ToggleServiceSpesificationActiveAndReturn(ctx, serviceID, specID)
	if err != nil {

		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"specID":    specID,
			"error":     err.Error(),
		})
		return nil, err
	}

	return serviceSpec, nil
}

func (u *JasaUsecase) ToggleServiceSpesificationRequired(ctx context.Context, serviceID, specID int64) (*models.ServiceSpesification, error) {

	service, err := u.JasaService.GetServiceByIDFromDB(ctx, serviceID)
	if err != nil {
		log.Logger.Error("jasaUsecase: GetServiceByIDFromDB")
		return nil, err
	}

	if service == nil {
		return nil, errors.New("service tidak ditemukan")
	}

	serviceSpec, err := u.JasaService.GetServiceSpesificationByID(ctx, specID)
	if err != nil {
		log.Logger.Error("jasaUsecase: u.JasaService.GetServiceSpesificationByID")
		return nil, err
	}

	if serviceSpec == nil {
		return nil, errors.New("service spesification tidak ditemukan")
	}

	serviceSpec, err = u.JasaService.ToggleServiceSpesificationRequiredAndReturn(ctx, serviceID, specID)
	if err != nil {

		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"specID":    specID,
			"error":     err.Error(),
		})
		return nil, err
	}

	return serviceSpec, nil
}

func (u *JasaUsecase) GetServiceByIDFromRead(ctx context.Context, serviceID int64) (*models.Service, error) {
	service, err := u.JasaService.GetServiceByIDFromRead(ctx, serviceID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"serviceID": serviceID,
		})

		return nil, err
	}

	return service, nil
}
