package puskesmasrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPuskesmasRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IPuskesmasRepository
	CreatePuskesmas(ctx context.Context, data *models.Puskesmas) error
	GetAllPuskesmas(ctx context.Context, limit, offset int, search string, kecamatanId uuid.UUID) ([]models.PuskesmasWithTotal, int, error)
	GetPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.PuskesmasWithTotal, error)
	UpdatePuskesmas(ctx context.Context, puskesmasId uuid.UUID, updates map[string]interface{}) error
	DeletePuskesmas(ctx context.Context, puskesmasId uuid.UUID) error

	GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error)
	GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error)
}
