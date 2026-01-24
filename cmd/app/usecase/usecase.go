package usecase

import (
	"jasamc/cmd/app/service"
	"jasamc/cmd/app/storage"
)

type JasaUsecase struct {
	JasaService    service.JasaService
	StorageService storage.Storage
}

func NewJasaUsecase(
	jasaService service.JasaService,
	storageService storage.Storage,
) *JasaUsecase {
	return &JasaUsecase{
		JasaService:    jasaService,
		StorageService: storageService,
	}
}
