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
func (k *KecamatanRepositoryImpl) GetAllKecamatan(
	ctx context.Context,
	limit int,
	offset int,
	search string,
) ([]models.KecamatanWithTotal, int, error) {

	var (
		kecamatanList []models.KecamatanWithTotal
		totalData     int64
	)

	if limit <= 0 {
		limit = 10
	}

	// ============================
	// 1. COUNT Kecamatan saja (FAST)
	// ============================
	countQuery := k.db.WithContext(ctx).
		Model(&models.Kecamatan{})

	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery = countQuery.Where(
			"nama_kecamatan ILIKE ? OR kode_wilayah ILIKE ?",
			searchPattern,
			searchPattern,
		)
	}

	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// ============================
	// 2. DATA Query + Total Kelurahan
	// ============================
	dataQuery := k.db.WithContext(ctx).
		Table("kecamatan").
		Select(`
			kecamatan.*,
			(SELECT COUNT(1) FROM kelurahan WHERE kelurahan.kecamatan_id = kecamatan.id) as total_kelurahan,
    		(SELECT COUNT(1) FROM puskesmas WHERE puskesmas.kecamatan_id = kecamatan.id) as total_puskesmas
		`).
		Order("kecamatan.created_at DESC").
		Limit(limit).
		Offset(offset)

	if search != "" {
		searchPattern := "%" + search + "%"
		dataQuery = dataQuery.Where(
			"kecamatan.nama_kecamatan ILIKE ? OR kecamatan.kode_wilayah ILIKE ?",
			searchPattern,
			searchPattern,
		)
	}

	if err := dataQuery.Find(&kecamatanList).Error; err != nil {
		return nil, 0, err
	}

	return kecamatanList, int(totalData), nil
}

// GetKecamatanById implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) GetKecamatanById(ctx context.Context, kecamatanId uuid.UUID) (*models.KecamatanWithTotal, error) {
	var result models.KecamatanWithTotal

	err := k.db.WithContext(ctx).
		Model(&models.Kecamatan{}).
		Select(`
        kecamatan.*,
        (SELECT COUNT(1) FROM kelurahan WHERE kelurahan.kecamatan_id = kecamatan.id) as total_kelurahan,
    	(SELECT COUNT(1) FROM puskesmas WHERE puskesmas.kecamatan_id = kecamatan.id) as total_puskesmas
		`).
		Where("kecamatan.id = ?", kecamatanId).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteKecamatan implements [IKecamatanRepository].
func (k *KecamatanRepositoryImpl) DeleteKecamatan(ctx context.Context, kecamatanId uuid.UUID) error {
	return k.db.WithContext(ctx).Where("id = ?", kecamatanId).Delete(&models.Kecamatan{}).Error
}
