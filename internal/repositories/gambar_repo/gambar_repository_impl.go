package gambarrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GambarRepositoryImpl struct {
	db *gorm.DB
}

func NewGambarRepositoryImpl(db *gorm.DB) IGambarRepository {
	return &GambarRepositoryImpl{
		db: db,
	}
}

// WithTx implements [IGambarRepository].
func (g *GambarRepositoryImpl) WithTx(tx *gorm.DB) IGambarRepository {
	return NewGambarRepositoryImpl(tx)
}

// DeleteByIds implements [IGambarRepository].
func (g *GambarRepositoryImpl) DeleteByIds(ctx context.Context, ids []uuid.UUID) error {
	return g.db.WithContext(ctx).Where("id IN ?", ids).Delete(&models.Gambar{}).Error
}
