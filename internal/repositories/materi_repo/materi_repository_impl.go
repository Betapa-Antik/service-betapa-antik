package materirepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MateriRepositoryImpl struct {
	db *gorm.DB
}

func NewMateriRepositoryImpl(db *gorm.DB) IMateriRepository {
	return &MateriRepositoryImpl{
		db: db,
	}
}

// DB implements [IMateriRepository].
func (m *MateriRepositoryImpl) DB() *gorm.DB {
	return m.db
}

// WithTx implements [IMateriRepository].
func (m *MateriRepositoryImpl) WithTx(tx *gorm.DB) IMateriRepository {
	return NewMateriRepositoryImpl(tx)
}

// CreateMateri implements [IMateriRepository].
func (m *MateriRepositoryImpl) CreateMateri(ctx context.Context, data *models.Materi) error {
	data.ID = uuid.New()
	return m.db.WithContext(ctx).Create(data).Error
}

// CreateGambar implements [IMateriRepository].
func (m *MateriRepositoryImpl) CreateGambar(ctx context.Context, data *models.Gambar) error {
	data.ID = uuid.New()
	return m.db.WithContext(ctx).Create(data).Error
}

// CreatePivotMateriGambar implements [IMateriRepository].
func (m *MateriRepositoryImpl) CreatePivotMateriGambar(ctx context.Context, data *models.MateriGambar) error {
	data.ID = uuid.New()
	return m.db.WithContext(ctx).Create(data).Error
}

// GetAllMateri implements [IMateriRepository].
func (m *MateriRepositoryImpl) GetAllMateri(ctx context.Context, limit, offset int, search string) ([]*models.Materi, int, error) {
	var (
		materiList []*models.Materi
		count      int64
	)

	if limit <= 0 {
		limit = 10
	}

	query := m.db.WithContext(ctx).
		Model(&models.Materi{}).
		Preload("Gambar")

	if search != "" {
		query = query.Where("judul ILIKE ?", "%"+search+"%")
	}

	// count pakai query yg sama
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// find pakai query yg sama
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&materiList).Error; err != nil {
		return nil, 0, err
	}

	return materiList, int(count), nil

}

// GetByID implements [IMateriRepository].
func (m *MateriRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Materi, error) {
	var materi models.Materi
	err := m.db.WithContext(ctx).Preload("Gambar").Preload("MateriGambars").Where("id = ?", id).First(&materi).Error
	if err != nil {
		return nil, err
	}
	return &materi, nil
}

// UpdateMateri implements [IMateriRepository].
func (m *MateriRepositoryImpl) UpdateMateri(ctx context.Context, id uuid.UUID, data *models.Materi) error {
	return m.db.WithContext(ctx).Model(&models.Materi{}).Where("id = ?", id).Updates(data).Error
}

// UpdateStatusMateri implements [IMateriRepository].
func (m *MateriRepositoryImpl) UpdateStatusMateri(ctx context.Context, id uuid.UUID, status string) error {
	return m.db.WithContext(ctx).Model(&models.Materi{}).Where("id = ?", id).Update("status", status).Error
}

// DeleteMateri implements [IMateriRepository].
func (m *MateriRepositoryImpl) DeleteMateri(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Materi{}).Error
}
