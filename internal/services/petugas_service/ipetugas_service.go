package petugasservice

import (
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	"betapa-antik-service/internal/models"
	"context"
)

type IPetugasService interface {
	GetSelectPuskesmas(ctx context.Context, search string) ([]models.SelectPuskesmas, error)
	RegisterAkunPetugas(ctx context.Context, req petugasrequest.RegisterPetugasRequest) error
}
