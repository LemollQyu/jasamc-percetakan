package service

import (
	"context"
	"jasamc/infrastructure/log"
	"jasamc/models"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *JasaService) GetAllServices(
	ctx context.Context,
) ([]models.Service, error) {

	// 1. Redis GET
	services, err := s.JasaRepo.GetAllServicesFromRedis(ctx)
	if err != nil {
		log.Logger.Error("REDIS GET ERROR (services:all)", err)
	}

	if len(services) > 0 {
		log.Logger.Info("CACHE HIT services:all")
		return services, nil
	}

	// 2. DB
	services, err = s.JasaRepo.GetAllServices(ctx)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("jasaService: gagal di s.JasaRepo.GetAllServices")
		return nil, err
	}

	log.Logger.Infof("CACHE MISS services:all → DB RETURNED %d ROWS", len(services))

	// 3. Redis SET ASYNC (DETACHED)
	go func(data []models.Service) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.JasaRepo.SetAllServicesToRedis(bgCtx, data); err != nil {
			log.Logger.Error("REDIS SET ERROR services:all (BG)", err)
		} else {
			log.Logger.Info("REDIS SET services:all SUCCESS (BG)")
		}
	}(services)

	return services, nil
}

func (s *JasaService) GetServiceByName(ctx context.Context, name string) (*models.Service, error) {
	jasa, err := s.JasaRepo.FindServiceByName(ctx, name)
	if err != nil {
		log.Logger.Error("jasaService: gagal di s.JasaRepo.FindJasaByName")
		return nil, err
	}

	return jasa, nil
}

func (s *JasaService) GetServiceBySlug(ctx context.Context, name string) (*models.Service, error) {
	jasa, err := s.JasaRepo.FindServiceBySlug(ctx, name)
	if err != nil {
		log.Logger.Error("jasaService: gagal di s.JasaRepo.FindServiceJasaBySlug")
		return nil, err
	}

	return jasa, nil
}

// unutk service redis get service by ID
func (s *JasaService) GetServiceByIDFromRead(ctx context.Context, serviceID int64) (*models.Service, error) {
	// GET dari redis dulu
	service, err := s.JasaRepo.GetServiceByIDFromRedis(ctx, serviceID)
	if err != nil {
		log.Logger.Error("REDIS GET ERROR (by id)", err)
	}

	if service != nil {
		log.Logger.WithFields(logrus.Fields{
			"id": serviceID,
		}).Info("CACHE HIT service:id")
		return service, nil
	}

	jasa, err := s.JasaRepo.FindServiceByID(ctx, serviceID)
	if err != nil {
		log.Logger.Error("jasaService: gagal di s.JasaRepo.FindServiceByID")
		return nil, err
	}

	if jasa == nil {
		return nil, nil
	}
	log.Logger.WithField("id", service).Infof("CACHE MISS service:%d", serviceID)

	// set redis async
	go func(data *models.Service) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.JasaRepo.SetServiceByIDToRedis(bgCtx, data)

		if err != nil {
			log.Logger.Errorf("REDIS SET ERROR by(%d)", serviceID)
		} else {

			log.Logger.Infof("REDIS SET services:%d SUCCESS (BG)", serviceID)

		}
	}(jasa)

	return jasa, nil
}

// get service by id ke database
func (s *JasaService) GetServiceByIDFromDB(ctx context.Context, serviceID int64) (*models.Service, error) {
	jasa, err := s.JasaRepo.FindServiceByID(ctx, serviceID)
	if err != nil {
		log.Logger.Error("jasaService: gagal di s.JasaRepo.FindServiceByID")
		return nil, err
	}

	return jasa, nil
}

func (s *JasaService) GetServiceMediaByID(ctx context.Context, mediaID int64) (*models.ServiceMedia, error) {
	media, err := s.JasaRepo.FindMediaServiceByID(ctx, mediaID)
	if err != nil {
		log.Logger.Error("jasaService: gagal di s.JasaRepo.FindMediaServiceByID")
		return nil, err
	}

	return media, nil
}

func (s *JasaService) CreateService(
	ctx context.Context,
	service *models.Service,
) error {
	err := s.JasaRepo.CreateServiceWithMedia(ctx, service)
	if err != nil {
		log.Logger.Error("JasaService: s.JasaRepo.CreateServiceWithMedia")
		return err
	}

	// delete redisnya
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, service.ID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", service.ID)
	}

	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return nil
}

func (s *JasaService) DeleteMediaBYServiceIDandMediaID(ctx context.Context, serviceID, mediaID int64) error {
	err := s.JasaRepo.DeleteMediaBYServiceIDandMediaID(ctx, serviceID, mediaID)
	if err != nil {
		return err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redis
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return nil
}

func (s *JasaService) UpdateStatusService(
	ctx context.Context,
	serviceID int64,
) (*models.Service, error) {

	service, err := s.JasaRepo.ToggleServiceActiveAndReturn(ctx, serviceID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"service_id": serviceID,
			"error":      err.Error(),
		}).Error("gagal di service ToggleJasaActive")

		return nil, err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all services in redis
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return service, nil
}

func (s *JasaService) CreateServiceSpesification(ctx context.Context, param *models.RequestServiceSpesification) (int64, error) {
	specID, err := s.JasaRepo.InsertServiceSpesification(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
			"error": err.Error(),
		})

		return 0, err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, param.ServiceId)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", param.ServiceId)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return specID, nil
}

func (s *JasaService) GetServiceSpesificationByName(ctx context.Context, serviceID int64, name string) (*models.ServiceSpesification, error) {
	serviceSpec, err := s.JasaRepo.FindServiceSpesificationByName(ctx, serviceID, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"name":      name,
			"error":     err.Error(),
		})
		return nil, err
	}

	return serviceSpec, nil

}

func (s *JasaService) AddServiceMedia(ctx context.Context, serviceID int64, param models.RequestAddServiceMedia) error {
	err := s.JasaRepo.AddNewServiceMedia(ctx, serviceID, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"param":     param,
			"error":     err.Error(),
		})
		return err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return nil
}

func (s *JasaService) CountServiceMediaByType(ctx context.Context, serviceID int64, typeMedia string) (int64, error) {
	count, err := s.JasaRepo.CountServiceMediaByType(ctx, serviceID, typeMedia)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"param":     typeMedia,
			"error":     err.Error(),
		})
		return 0, nil
	}

	return count, nil
}

func (s *JasaService) CreateServiceSpesificationValue(ctx context.Context, param models.RequestServiceSpesificationValue) error {
	err := s.JasaRepo.InsertNewServiceSpesificationValue(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
			"error": err.Error(),
		})
		return err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, param.ServiceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", param.ServiceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return nil
}

func (s *JasaService) GetServiceSpesificationByID(ctx context.Context, serviceSpecID int64) (*models.ServiceSpesification, error) {
	serviceSpec, err := s.JasaRepo.FindServiceSpesificationByID(ctx, serviceSpecID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"id":    serviceSpecID,
			"error": err.Error(),
		})
		return nil, err
	}

	return serviceSpec, nil
}

func (s *JasaService) GetServiceSpesificationValueByName(ctx context.Context, serviceID, spesificationID int64, value string) (*models.ServiceSpesificationValue, error) {

	serviceSpec, err := s.JasaRepo.FindServiceSpesificationValueByServiceAndSpec(ctx, serviceID, spesificationID, value)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"ServiceID":       serviceID,
			"SpesificationID": spesificationID,
			"Value":           value,
			"error":           err.Error(),
		})
	}
	return serviceSpec, err
}

func (s *JasaService) GetServiceSpesificationValueByServiceIDBySpecIDByValueID(ctx context.Context, serviceID, specID, valueID int64) (*models.ServiceSpesificationValue, error) {
	specValue, err := s.JasaRepo.FindServiceSpesificationValueByServiceIDBySpecIDByValueID(ctx, serviceID, specID, valueID)

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"ServiceID":       serviceID,
			"SpesificationID": specID,
			"ValueID":         valueID,
			"error":           err.Error(),
		})

		return nil, err
	}

	return specValue, nil
}

func (s *JasaService) UpdateServiceSpesificationValue(
	ctx context.Context,
	serviceID int64,
	specID int64,
	valueID int64,
	param models.RequestUpdateServiceSpesificationValue,
) error {
	err := s.JasaRepo.UpdateServiceSpesificationValue(
		ctx,
		serviceID,
		specID,
		valueID,
		param,
	)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"serviceID": serviceID,
			"specID":    specID,
			"valueID":   valueID,
			"param":     param,
		})

		return err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return nil
}

func (s *JasaService) DeleteServiceSpesification(ctx context.Context, serviceID, specID int64) error {
	err := s.JasaRepo.DeleteServiceSpesification(ctx, serviceID, specID)
	if err != nil {
		log.Logger.Error("jasaService: s.JasaRepo.DeleteServiceSpesification")
		return err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return nil
}

func (s *JasaService) DeleteServiceByID(ctx context.Context, serviceID int64) error {
	err := s.JasaRepo.DeleteServiceByID(ctx, serviceID)
	if err != nil {
		log.Logger.Error("JasaService: s.JasaService.DeleteServiceByID")
		return err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}
	return nil
}

func (s *JasaService) ToggleServiceSpesificationActiveAndReturn(ctx context.Context, serviceID, specID int64) (*models.ServiceSpesification, error) {
	serviceSpec, err := s.JasaRepo.ToggleServiceSpesificationActiveAndReturn(ctx, serviceID, specID)
	if err != nil {

		log.Logger.Error("JasaService: s.JasaRepo.ToggleServiceSpesificationActiveAndReturn")
		return nil, err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return serviceSpec, nil
}

func (s *JasaService) ToggleServiceSpesificationRequiredAndReturn(ctx context.Context, serviceID, specID int64) (*models.ServiceSpesification, error) {
	serviceSpec, err := s.JasaRepo.ToggleServiceSpesificationRequiredAndReturn(ctx, serviceID, specID)
	if err != nil {

		log.Logger.Error("JasaService: s.JasaRepo.ToggleServiceSpesificationRequiredAndReturn")
		return nil, err
	}

	// delete data service id redis
	err = s.JasaRepo.DeleteServiceCacheByID(ctx, serviceID)
	if err != nil {
		log.Logger.Errorf("gagal hapus redis service:%d", serviceID)
	}

	// delete data all redia
	if err := s.JasaRepo.DeleteAllServiceCache(ctx); err != nil {
		log.Logger.Error("Gagal menghapus cache Redis services:all", err)
	}

	return serviceSpec, nil
}
