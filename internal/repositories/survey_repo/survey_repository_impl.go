package surveyrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SurveyRepositoryImpl struct {
	db *gorm.DB
}

func NewSurveyRepositoryImpl(db *gorm.DB) ISurveyRepository {
	return &SurveyRepositoryImpl{db: db}
}

// DB implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) DB() *gorm.DB {
	return s.db
}

// WithTx implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) WithTx(tx *gorm.DB) ISurveyRepository {
	return NewSurveyRepositoryImpl(tx)
}

// GetSelectKeluarga implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) GetSelectKeluarga(ctx context.Context, search string) ([]models.SelectKeluarga, error) {
	var result []models.SelectKeluarga

	query := s.db.WithContext(ctx).
		Table("keluarga").
		Select(`
			keluarga.id,
			keluarga.nama_kepala_keluarga,
			kecamatan.nama_kecamatan as kecamatan,
			kelurahan.nama_kelurahan as kelurahan,
			keluarga.rt,
			keluarga.rw,
			keluarga.alamat
		`).
		Joins("LEFT JOIN kecamatan ON kecamatan.id = keluarga.kecamatan_id").
		Joins("LEFT JOIN kelurahan ON kelurahan.id = keluarga.kelurahan_id").
		Order("keluarga.nama_kepala_keluarga ASC")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("keluarga.nama_kepala_keluarga ILIKE ?", searchPattern)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// CreateSurvey implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) CreateSurvey(ctx context.Context, data *models.Survey) error {
	data.ID = uuid.New()
	return s.db.WithContext(ctx).Create(data).Error
}

// CreateSurveyItem implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) CreateSurveyItem(ctx context.Context, data *models.SurveyItem) error {
	data.ID = uuid.New()
	return s.db.WithContext(ctx).Create(data).Error
}

// // FindSurveyLokasi implements [ISurveyRepository].
// func (s *SurveyRepositoryImpl) FindSurveyLokasi(ctx context.Context, NamaLokasi string, JenisSurvey string) (*models.SurveyLokasi, error) {
// 	var lokasi models.SurveyLokasi
// 	if err := s.db.WithContext(ctx).Where("nama_lokasi = ? AND jenis_survey = ?", NamaLokasi, JenisSurvey).First(&lokasi).Error; err != nil {
// 		return nil, err
// 	}

// 	return &lokasi, nil
// }

// // CreateSurveyLokasi implements [ISurveyRepository].
// func (s *SurveyRepositoryImpl) CreateSurveyLokasi(ctx context.Context, data *models.SurveyLokasi) error {
// 	data.ID = uuid.New()
// 	return s.db.WithContext(ctx).Create(data).Error
// }

// FindOrCreateSurveyLokasi implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) FindOrCreateSurveyLokasi(ctx context.Context, namaLokasi string, jenisSurvey string) (*models.SurveyLokasi, error) {
	var lokasi models.SurveyLokasi

	err := s.db.WithContext(ctx).
		Where("nama_lokasi = ? AND jenis_survey = ?", namaLokasi, jenisSurvey).
		First(&lokasi).Error

	// Kalau ketemu → return
	if err == nil {
		return &lokasi, nil
	}

	// Kalau error selain not found
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Kalau belum ada → buat baru
	lokasi = models.SurveyLokasi{
		ID:          uuid.New(),
		NamaLokasi:  namaLokasi,
		JenisSurvey: jenisSurvey,
	}

	if err := s.db.WithContext(ctx).Create(&lokasi).Error; err != nil {
		return nil, err
	}

	return &lokasi, nil

}

// CreateSurveyFollowUpJentik implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) CreateSurveyFollowUpJentik(ctx context.Context, data *models.SurveyFollowUpJentik) error {
	data.ID = uuid.New()
	return s.db.WithContext(ctx).Create(data).Error
}

// CreateSurveyFollowUpNyamuk implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) CreateSurveyFollowUpNyamuk(ctx context.Context, data *models.SurveyFollowUpNyamuk) error {
	data.ID = uuid.New()
	return s.db.WithContext(ctx).Create(data).Error
}

// CreateGambar implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) CreateGambar(ctx context.Context, data *models.Gambar) error {
	data.ID = uuid.New()
	return s.db.WithContext(ctx).Create(data).Error
}

// CreatePivotSurveyGambar implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) CreatePivotSurveyGambar(ctx context.Context, data *models.SurveyGambar) error {
	data.ID = uuid.New()
	return s.db.WithContext(ctx).Create(data).Error
}

// GetAllSurvey implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) GetAllSurvey(ctx context.Context, limit int, offset int, search string, jenisSurvey string, startDate string, endDate string, petugasId uuid.UUID) ([]models.Survey, int, error) {
	var (
		surveys []models.Survey
		total   int64
	)

	if limit <= 0 {
		limit = 10
	}
	// 1. Gunakan Model di awal agar GORM tahu skema struct-nya
	query := s.db.WithContext(ctx).Model(&models.Survey{}).Where("petugas_id = ?", petugasId)

	// 2. Join keluarga untuk pencarian
	// Gunakan Alias jika perlu, tapi pastikan kolom yang di-select adalah milik survey.*
	query = query.Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id")

	if search != "" {
		searchLike := "%" + search + "%"
		// Gunakan scope table 'keluarga' secara eksplisit
		query = query.Where(`
            keluarga.nama_kepala_keluarga ILIKE ? OR
            keluarga.alamat ILIKE ? OR
            keluarga.kecamatan ILIKE ? OR
            keluarga.kelurahan ILIKE ?
        `, searchLike, searchLike, searchLike, searchLike)
	}

	if jenisSurvey != "" {
		query = query.Where("survey.jenis_survey = ?", jenisSurvey)
	}

	if startDate != "" && endDate != "" {
		query = query.Where("survey.tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	// Hitung total sebelum limit/offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. Ambil data dengan Preload
	// Tambahkan Select("survey.*") untuk menghindari tabrakan kolom ID dari hasil Join
	err := query.Select("survey.*").
		Preload("Keluarga").
		Preload("Keluarga.Kecamatan").
		Preload("Keluarga.Kelurahan").
		Preload("Petugas").
		Preload("Items.Lokasi").
		Preload("SurveyGambar.Gambar").
		Preload("FollowUpNyamuk").
		Preload("FollowUpJentik").
		Order("survey.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&surveys).Error

	if err != nil {
		return nil, 0, err
	}

	return surveys, int(total), nil
}

// GetSurveyByID implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) GetSurveyByID(ctx context.Context, surveyId uuid.UUID) (*models.Survey, error) {
	var result models.Survey
	if err := s.db.WithContext(ctx).
		Preload("Keluarga").
		Preload("Keluarga.Kecamatan").
		Preload("Keluarga.Kelurahan").
		Preload("Petugas").
		Preload("Items.Lokasi").
		Preload("SurveyGambar.Gambar").
		Preload("FollowUpNyamuk").
		Preload("FollowUpJentik").First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSurvey implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) UpdateSurvey(ctx context.Context, surveyId uuid.UUID, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&models.Survey{}).Where("id = ?", surveyId).Updates(updates).Error
}

// DeleteSurvey implements [ISurveyRepository].
func (s *SurveyRepositoryImpl) DeleteSurvey(ctx context.Context, surveyId uuid.UUID) error {
	return s.db.WithContext(ctx).Where("id = ?", surveyId).Delete(&models.Survey{}).Error
}
