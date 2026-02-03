package materirepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMateriRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IMateriRepository
	CreateMateri(ctx context.Context, data *models.Materi) error
	CreateGambar(ctx context.Context, data *models.Gambar) error
	CreatePivotMateriGambar(ctx context.Context, data *models.MateriGambar) error

	GetAllMateri(ctx context.Context, limit, offset int, search string) ([]*models.Materi, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Materi, error)
	UpdateMateri(ctx context.Context, id uuid.UUID, data *models.Materi) error
	UpdateStatusMateri(ctx context.Context, id uuid.UUID, status string) error
	DeleteMateri(ctx context.Context, id uuid.UUID) error
}
