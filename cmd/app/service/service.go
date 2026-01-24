package service

import "jasamc/cmd/app/repository"

type JasaService struct {
	JasaRepo repository.JasaRepository
}

func NewJasaService(jasaRepo repository.JasaRepository) *JasaService {
	return &JasaService{
		JasaRepo: jasaRepo,
	}
}
