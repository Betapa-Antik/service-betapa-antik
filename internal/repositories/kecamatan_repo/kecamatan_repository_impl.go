package kecamatanrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KecamatanRepositoryImpl struct {
	db *gorm.DB
}

func NewKecamatanRepositoryImpl(db *gorm.DB) IKecamatanRepository {
	return &KecamatanRepositoryImpl{db: db}
}

// DB implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) DB() *gorm.DB {
	return k.db
}

// WithTx implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) WithTx(tx *gorm.DB) IKecamatanRepository {
	return NewKecamatanRepositoryImpl(tx)
}

// CreateKecamatan implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) CreateKecamatan(ctx context.Context, data *models.Kecamatan) error {
	data.ID = uuid.New()
	return k.db.WithContext(ctx).Create(data).Error
}

// Update implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) Update(ctx context.Context, kecamatanId uuid.UUID, updates map[string]interface{}) error {
	return k.db.WithContext(ctx).Model(&models.Kecamatan{}).Where("id = ?", kecamatanId).Updates(updates).Error
}

// GetAllKecamatan implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) GetAllKecamatan(ctx context.Context, limit int, ofset int, search string) ([]*models.Kecamatan, int, error) {
	var (
		kecamatanList []*models.Kecamatan
		count         int64
	)

	if limit <= 0 {
		limit = 10
	}

	query := k.db.WithContext(ctx).Model(&models.Kecamatan{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama_kecamatan ILIKE ? OR kode_wilayah ILIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(ofset).Order("created_at DESC").Find(&kecamatanList).Error
	if err != nil {
		return nil, 0, err
	}
	return kecamatanList, int(count), nil
}

// GetKecamatanById implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) GetKecamatanById(ctx context.Context, kecamatanId uuid.UUID) (*models.Kecamatan, error) {
	var kecamatan models.Kecamatan
	err := k.db.WithContext(ctx).Where("id = ?", kecamatanId).First(&kecamatan).Error
	if err != nil {
		return nil, err
	}
	return &kecamatan, nil
}

// DeleteKecamatan implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) DeleteKecamatan(ctx context.Context, kecamatanId uuid.UUID) error {
	return k.db.WithContext(ctx).Where("id = ?", kecamatanId).Delete(&models.Kecamatan{}).Error
}
