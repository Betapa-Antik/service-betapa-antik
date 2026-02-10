package keluargaservice

import (
	keluargarequest "betapa-antik-service/internal/dto/request/keluarga_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IKeluargaService interface {
	CreateKeluarga(ctx context.Context, req keluargarequest.CreateKeluargaRequest) error
	GetAllKeluarga(ctx context.Context, req keluargarequest.GetAllKeluargaRequest) ([]models.Keluarga, int, error)
	GetKeluargaById(ctx context.Context, keluargaId uuid.UUID) (*models.Keluarga, error)
	UpdateKeluarga(ctx context.Context, keluargaId uuid.UUID, req keluargarequest.UpdateKeluargaRequest) error
	DeleteKeluarga(ctx context.Context, keluargaId uuid.UUID) error

	GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error)
	GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error)
}
