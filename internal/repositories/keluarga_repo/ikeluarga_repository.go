package keluargarepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IKeluargaRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IKeluargaRepository
	CreateKeluarga(ctx context.Context, data *models.Keluarga) error
	GetAllKeluarga(ctx context.Context, limit, offset int, search string) ([]models.Keluarga, int, error)
	GetKeluargaById(ctx context.Context, keluargaId uuid.UUID) (*models.Keluarga, error)
	UpdateKeluarga(ctx context.Context, keluargaId uuid.UUID, updates map[string]interface{}) error
	DeleteKeluarga(ctx context.Context, keluargaId uuid.UUID) error

	GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error)
	GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error)
}
