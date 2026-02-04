package kelurahanrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KelurahanRepositoryImpl struct {
	db *gorm.DB
}

func NewKelurahanRepositoryImpl(db *gorm.DB) IKelurahanRepository {
	return &KelurahanRepositoryImpl{db: db}
}

// DB implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) DB() *gorm.DB {
	return k.db
}

// WithTx implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) WithTx(tx *gorm.DB) IKelurahanRepository {
	return NewKelurahanRepositoryImpl(tx)
}

// CreateKelurahan implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) CreateKelurahan(ctx context.Context, data *models.Kelurahan) error {
	data.ID = uuid.New()
	return k.db.WithContext(ctx).Create(data).Error
}

// GetAllKelurahan implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) GetAllKelurahan(ctx context.Context, kecamatanId uuid.UUID, limit int, offset int, search string) ([]models.Kelurahan, int, error) {
	var (
		kelurahanList []models.Kelurahan
		count         int64
	)

	if limit <= 0 {
		limit = 10
	}
	query := k.db.WithContext(ctx).Model(&models.Kelurahan{}).Preload("Kecamatan").Where("kecamatan_id = ?", kecamatanId)

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama_kelurahan ILIKE ? OR kode_wilayah ILIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&kelurahanList).Error
	if err != nil {
		return nil, 0, err
	}

	return kelurahanList, int(count), nil
}

// GetKelurahanById implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) GetKelurahanById(ctx context.Context, kelurahanId uuid.UUID) (*models.Kelurahan, error) {
	var kelurahan models.Kelurahan
	err := k.db.WithContext(ctx).Preload("Kecamatan").Where("id = ?", kelurahanId).First(&kelurahan).Error
	if err != nil {
		return nil, err
	}
	return &kelurahan, nil
}

// UpdateKelurahan implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) UpdateKelurahan(ctx context.Context, kelurahanId uuid.UUID, updates map[string]interface{}) error {
	return k.db.WithContext(ctx).Model(&models.Kelurahan{}).Where("id = ?", kelurahanId).Updates(updates).Error
}

// DeleteKelurahan implements [IKelurahanRepository].
func (k *KelurahanRepositoryImpl) DeleteKelurahan(ctx context.Context, kelurahanId uuid.UUID) error {
	return k.db.WithContext(ctx).Where("id = ?", kelurahanId).Delete(&models.Kelurahan{}).Error
}
