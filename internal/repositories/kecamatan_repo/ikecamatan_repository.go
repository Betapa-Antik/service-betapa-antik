package kecamatanrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IKecamatanRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IKecamatanRepository
	CreateKecamatan(ctx context.Context, data *models.Kecamatan) error
	GetAllKecamatan(ctx context.Context, limit, ofset int, search string) ([]models.KecamatanWithTotal, int, error)
	GetKecamatanById(ctx context.Context, kecamatanId uuid.UUID) (*models.KecamatanWithTotal, error)
	Update(ctx context.Context, kecamatanId uuid.UUID, updates map[string]interface{}) error
	DeleteKecamatan(ctx context.Context, kecamatanId uuid.UUID) error
}
