package handler

import "jasamc/cmd/app/usecase"

type JasaHandler struct {
	JasaUsecase usecase.JasaUsecase
}

func NewJasaHandler(jasaUsecase usecase.JasaUsecase) *JasaHandler {
	return &JasaHandler{
		JasaUsecase: jasaUsecase,
	}
}
