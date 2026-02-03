package materigambarrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMateriGambarRepository interface {
	WithTx(tx *gorm.DB) IMateriGambarRepository
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]models.MateriGambar, error)
	DeleteByIds(ctx context.Context, ids []uuid.UUID) error
}
