package materigambarrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MateriGambarRepositoryImpl struct {
	db *gorm.DB
}

func NewMateriGambarRepositoryImpl(db *gorm.DB) IMateriGambarRepository {
	return &MateriGambarRepositoryImpl{
		db: db,
	}
}

// WithTx implements [IMateriGambarRepository].
func (m *MateriGambarRepositoryImpl) WithTx(tx *gorm.DB) IMateriGambarRepository {
	return NewMateriGambarRepositoryImpl(tx)
}

// FindByIds implements [IMateriGambarRepository].
func (m *MateriGambarRepositoryImpl) FindByIds(ctx context.Context, ids []uuid.UUID) ([]models.MateriGambar, error) {
	var materiGambar []models.MateriGambar
	err := m.db.WithContext(ctx).
		Select("id", "gambar_id").
		Preload("Gambar", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "path")
		}).
		Where("id IN ?", ids).
		Find(&materiGambar).Error

	if err != nil {
		return nil, err
	}
	return materiGambar, nil
}

// DeleteByIds implements [IMateriGambarRepository].
func (m *MateriGambarRepositoryImpl) DeleteByIds(ctx context.Context, ids []uuid.UUID) error {
	return m.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.MateriGambar{}).Error
}
