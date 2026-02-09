package puskesmasrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PuskesmasRepositoryImpl struct {
	db *gorm.DB
}

func NewPuskesmasRepositoryImpl(db *gorm.DB) IPuskesmasRepository {
	return &PuskesmasRepositoryImpl{db: db}
}

// DB implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) DB() *gorm.DB {
	return p.db
}

// WithTx implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) WithTx(tx *gorm.DB) IPuskesmasRepository {
	return NewPuskesmasRepositoryImpl(tx)
}

// CreatePuskesmas implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) CreatePuskesmas(ctx context.Context, data *models.Puskesmas) error {
	data.ID = uuid.New()
	return p.db.WithContext(ctx).Create(data).Error
}

// UpdatePuskesmas implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) UpdatePuskesmas(ctx context.Context, puskesmasId uuid.UUID, updates map[string]interface{}) error {
	return p.db.WithContext(ctx).Model(&models.Puskesmas{}).Where("id = ?", puskesmasId).Updates(updates).Error
}

// GetAllPuskesmas implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) GetAllPuskesmas(ctx context.Context, limit int, offset int, search string, kecamatanId uuid.UUID) ([]models.PuskesmasWithTotal, int, error) {
	var (
		puskesmasList []models.PuskesmasWithTotal
		totalData     int64
	)

	if limit <= 0 {
		limit = 10
	}

	baseQuery := p.db.WithContext(ctx).
		Model(&models.Puskesmas{}).Preload("Kecamatan").Preload("Kelurahan")

	if search != "" {
		searchPattern := "%" + search + "%"
		baseQuery = baseQuery.Where("puskesmas.nama_puskesmas ILIKE ?", searchPattern)
	}

	if kecamatanId != uuid.Nil {
		baseQuery = baseQuery.Where("puskesmas.kecamatan_id = ?", kecamatanId)
	}

	if err := baseQuery.
		Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := baseQuery.
		Select(`
			puskesmas.*,
			(
				SELECT COUNT(1)
				FROM "user"
				WHERE "user".puskesmas_id = puskesmas.id
			) as total_petugas
		`).
		Order("puskesmas.created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := dataQuery.Find(&puskesmasList).Error; err != nil {
		return nil, 0, err
	}

	return puskesmasList, int(totalData), nil
}

// GetPuskesmasById implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) GetPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.PuskesmasWithTotal, error) {
	var result models.PuskesmasWithTotal
	err := p.db.WithContext(ctx).Model(&models.Puskesmas{}).
		Preload("Kecamatan").Preload("Kelurahan").
		Select(`
			puskesmas.*,
			(SELECT COUNT(1) from "user" WHERE "user".puskesmas_id = puskesmas.id) as total_petugas
		`).
		Where("puskesmas.id = ?", puskesmasId).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeletePuskesmas implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) DeletePuskesmas(ctx context.Context, puskesmasId uuid.UUID) error {
	return p.db.WithContext(ctx).Where("id = ?", puskesmasId).Delete(&models.Puskesmas{}).Error
}

// GetSelectKecamatan implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error) {
	var result []models.SelectKecamatan

	query := p.db.WithContext(ctx).
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

// GetSelectKelurahan implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error) {
	var result []models.SelectKelurahan

	query := p.db.WithContext(ctx).
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
