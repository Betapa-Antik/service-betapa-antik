package keluargarepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeluargaRepositoryImpl struct {
	db *gorm.DB
}

func NewKeluargaRepositoryImpl(db *gorm.DB) IKeluargaRepository {
	return &KeluargaRepositoryImpl{db: db}
}

// DB implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) DB() *gorm.DB {
	return k.db
}

// WithTx implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) WithTx(tx *gorm.DB) IKeluargaRepository {
	return NewKeluargaRepositoryImpl(tx)
}

// CreateKeluarga implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) CreateKeluarga(ctx context.Context, data *models.Keluarga) error {
	data.ID = uuid.New()
	return k.db.WithContext(ctx).Create(data).Error
}

// GetAllKeluarga implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) GetAllKeluarga(ctx context.Context, limit int, offset int, search string) ([]models.Keluarga, int, error) {
	var (
		keluargaList []models.Keluarga
		count        int64
	)

	if limit <= 0 {
		limit = 10
	}

	query := k.db.WithContext(ctx).Model(&models.Keluarga{}).Preload("Kecamatan").Preload("Kelurahan")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama_kepala_keluarga ILIKE ?", searchPattern)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&keluargaList).Error
	if err != nil {
		return nil, 0, err
	}

	return keluargaList, int(count), nil
}

// GetKeluargaById implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) GetKeluargaById(ctx context.Context, keluargaId uuid.UUID) (*models.Keluarga, error) {
	var keluarga models.Keluarga
	err := k.db.WithContext(ctx).Preload("Kecamatan").Preload("Kelurahan").Where("id = ?", keluargaId).First(&keluarga).Error
	if err != nil {
		return nil, err
	}

	return &keluarga, nil
}

// UpdateKeluarga implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) UpdateKeluarga(ctx context.Context, keluargaId uuid.UUID, updates map[string]interface{}) error {
	return k.db.WithContext(ctx).Model(&models.Keluarga{}).Where("id = ?", keluargaId).Updates(updates).Error
}

// DeleteKeluarga implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) DeleteKeluarga(ctx context.Context, keluargaId uuid.UUID) error {
	return k.db.WithContext(ctx).Where("id = ?", keluargaId).Delete(&models.Keluarga{}).Error
}

// GetSelectKecamatan implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error) {
	var result []models.SelectKecamatan

	query := k.db.WithContext(ctx).
		Table("kecamatan").
		Select(`
			kecamatan.id,
			kecamatan.nama_kecamatan,
			kecamatan.kode_wilayah
		`).Order("kecamatan.nama_kecamatan ASC")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("kecamatan.nama_kecamatan ILIKE ?", searchPattern)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// GetSelectKelurahan implements [IKeluargaRepository].
func (k *KeluargaRepositoryImpl) GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error) {
	var result []models.SelectKelurahan

	query := k.db.WithContext(ctx).
		Table("kelurahan").
		Select(`
		kelurahan.id,
		kelurahan.nama_kelurahan,
		kelurahan.kode_kelurahan
	`).
		Where("kelurahan.kecamatan_id = ?", kecamatanId).
		Order("kelurahan.nama_kelurahan ASC")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("kelurahan.nama_kelurahan ILIKE ?", searchPattern)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
