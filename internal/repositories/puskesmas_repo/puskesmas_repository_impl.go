package puskesmasrepo

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"
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
func (p *PuskesmasRepositoryImpl) GetAllPuskesmas(
	ctx context.Context,
	limit int,
	offset int,
	search string,
	kecamatanId uuid.UUID,
) ([]models.PuskesmasWithTotal, int, error) {

	var (
		puskesmasList []models.PuskesmasWithTotal
		totalData     int64
	)

	if limit <= 0 {
		limit = 10
	}

	// ============================
	// BASE QUERY
	// ============================
	baseQuery := p.db.WithContext(ctx).
		Model(&models.Puskesmas{}).
		Preload("Kecamatan").
		Preload("Kelurahan")

	if search != "" {
		searchPattern := "%" + search + "%"
		baseQuery = baseQuery.Where("puskesmas.nama_puskesmas ILIKE ?", searchPattern)
	}

	if kecamatanId != uuid.Nil {
		baseQuery = baseQuery.Where("puskesmas.kecamatan_id = ?", kecamatanId)
	}

	// ============================
	// COUNT
	// ============================
	if err := baseQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// ============================
	// DATA QUERY + DF Kecamatan
	// ============================
	dataQuery := baseQuery.
		Joins("LEFT JOIN kecamatan ON kecamatan.id = puskesmas.kecamatan_id").
		Joins("LEFT JOIN kelurahan ON kelurahan.id = puskesmas.kelurahan_id").
		Joins("LEFT JOIN keluarga ON keluarga.kecamatan_id = puskesmas.kecamatan_id").
		Joins(`
        LEFT JOIN survey 
        ON survey.keluarga_id = keluarga.id
        AND survey.jenis_survey = ?
    `, models.JenisSurveyJentik).
		Joins("LEFT JOIN survey_item ON survey_item.survey_id = survey.id").
		Select(`
        puskesmas.*,

        kecamatan.nama_kecamatan,
        kelurahan.nama_kelurahan,

        -- Total Petugas Active
        (SELECT COUNT(1)
         FROM "user"
         WHERE "user".puskesmas_id = puskesmas.id
           AND "user".status = 'active'
        ) as total_petugas,

        -- DF Kecamatan (Max DF)
        COALESCE(
            GREATEST(
                CASE
                    WHEN (
                        COUNT(DISTINCT CASE 
                            WHEN survey_item.jumlah_positif > 0 
                            THEN survey.id
                        END)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 3 THEN 1
                    WHEN (
                        COUNT(DISTINCT CASE 
                            WHEN survey_item.jumlah_positif > 0 
                            THEN survey.id
                        END)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 7 THEN 2
                    WHEN (
                        COUNT(DISTINCT CASE 
                            WHEN survey_item.jumlah_positif > 0 
                            THEN survey.id
                        END)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 17 THEN 3
                    WHEN (
                        COUNT(DISTINCT CASE 
                            WHEN survey_item.jumlah_positif > 0 
                            THEN survey.id
                        END)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 28 THEN 4
                    ELSE 9
                END,

                CASE
                    WHEN (
                        SUM(survey_item.jumlah_positif)::float
                        / NULLIF(SUM(survey_item.jumlah_tempat_air),0) * 100
                    ) <= 2 THEN 1
                    WHEN (
                        SUM(survey_item.jumlah_positif)::float
                        / NULLIF(SUM(survey_item.jumlah_tempat_air),0) * 100
                    ) <= 5 THEN 2
                    WHEN (
                        SUM(survey_item.jumlah_positif)::float
                        / NULLIF(SUM(survey_item.jumlah_tempat_air),0) * 100
                    ) <= 9 THEN 3
                    ELSE 9
                END,

                CASE
                    WHEN (
                        SUM(survey_item.jumlah_positif)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 4 THEN 1
                    WHEN (
                        SUM(survey_item.jumlah_positif)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 9 THEN 2
                    WHEN (
                        SUM(survey_item.jumlah_positif)::float
                        / NULLIF(COUNT(DISTINCT survey.id),0) * 100
                    ) <= 19 THEN 3
                    ELSE 9
                END
            ),
        1) as df
    `).
		Group("puskesmas.id, kecamatan.nama_kecamatan, kelurahan.nama_kelurahan").
		Order("puskesmas.created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := dataQuery.Scan(&puskesmasList).Error; err != nil {
		return nil, 0, err
	}

	// ============================
	// STATUS dari DF
	// ============================
	for i := range puskesmasList {
		puskesmasList[i].Status = utils.GetStatusByDF(int(puskesmasList[i].DF))
	}

	return puskesmasList, int(totalData), nil
}

// GetPuskesmasById implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) GetPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.PuskesmasWithTotal, error) {
	var result models.PuskesmasWithTotal
	err := p.db.WithContext(ctx).Model(&models.Puskesmas{}).
		Preload("Kecamatan").
		Preload("Kelurahan").
		Select(`
        puskesmas.*,
        (SELECT COUNT(1)
         FROM "user"
         WHERE "user".puskesmas_id = puskesmas.id
           AND "user".status = 'active') as total_petugas
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

// GetAllPetugasByPuskesmasId implements [IPuskesmasRepository].
func (p *PuskesmasRepositoryImpl) GetAllPetugasByPuskesmasId(
	ctx context.Context,
	puskesmasId uuid.UUID,
	search string,
) ([]*models.PetugasWithTotalSurvey, error) {

	var petugasList []*models.PetugasWithTotalSurvey

	query := p.db.WithContext(ctx).
		Table(`"user"`).
		Joins(`LEFT JOIN role ON role.id = "user".role_id`).
		Joins(`LEFT JOIN puskesmas ON puskesmas.id = "user".puskesmas_id`).
		Joins(`LEFT JOIN survey ON survey.petugas_id = "user".id`).
		Where(`"user".puskesmas_id = ?`, puskesmasId).
		Where(`role.nama = ?`, "PETUGAS PUSKESMAS").
		Select(`
			"user".id,
			"user".foto,
			"user".nama_lengkap,
			"user".no_pegawai,
			"user".email,
			"user".status,
			"user".created_at,
			"user".updated_at,

			puskesmas.nama_puskesmas as puskesmas,
			role.nama as jabatan,

			COUNT(DISTINCT survey.id) as total_survey
		`).
		Group(`
			"user".id,
			puskesmas.nama_puskesmas,
			role.nama
		`).
		Order(`"user".created_at DESC`)

	// 🔍 Search Nama atau Email
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(`
			"user".nama_lengkap ILIKE ?
			OR "user".email ILIKE ?
		`, searchPattern, searchPattern)
	}

	if err := query.Scan(&petugasList).Error; err != nil {
		return nil, err
	}

	return petugasList, nil
}
