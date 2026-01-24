package grpc

import (
	"context"
	"jasamc/cmd/app/usecase"
	"jasamc/proto/jasapb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	jasapb.UnimplementedJasaServiceServer
	JasaUsecase usecase.JasaUsecase
}

func (s *GRPCServer) GetServiceDetail(
	ctx context.Context,
	req *jasapb.GetServiceDetailRequest,
) (*jasapb.GetServiceDetailResponse, error) {

	service, err := s.JasaUsecase.GetServiceByIDFromRead(ctx, req.ServiceId)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	specs := make([]*jasapb.Spesification, 0)

	for _, spec := range service.Spesification {

		values := make([]*jasapb.SpesificationValue, 0)

		for _, v := range spec.SpesificationValue {
			values = append(values, &jasapb.SpesificationValue{
				Id:              v.ID,
				Value:           v.Value,
				AdditionalPrice: int64(v.AdditionalPrice),
			})
		}

		specs = append(specs, &jasapb.Spesification{
			Id:         spec.ID,
			Name:       spec.Name,
			InputType:  spec.InputType,
			IsRequired: spec.IsRequired,
			Values:     values,
		})
	}

	return &jasapb.GetServiceDetailResponse{
		Id:             service.ID,
		Name:           service.Name,
		BasePrice:      int64(service.BasePrice),
		Spesifications: specs,
	}, nil
}
