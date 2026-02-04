package kelurahanrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IKelurahanRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IKelurahanRepository
	CreateKelurahan(ctx context.Context, data *models.Kelurahan) error
	GetAllKelurahan(ctx context.Context, kecamatanId uuid.UUID, limit, offset int, search string) ([]*models.Kelurahan, int, error)
	GetKelurahanById(ctx context.Context, kelurahanId uuid.UUID) (*models.Kelurahan, error)
	UpdateKelurahan(ctx context.Context, kelurahanId uuid.UUID, updates map[string]interface{}) error
	DeleteKelurahan(ctx context.Context, kelurahanId uuid.UUID) error
}
