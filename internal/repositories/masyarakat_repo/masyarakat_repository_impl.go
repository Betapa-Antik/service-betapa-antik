package masyarakatrepo

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasyarakatRepositoryImpl struct {
	db *gorm.DB
}

func NewMasyarakatRepositoryImpl(db *gorm.DB) IMasyarakatRepository {
	return &MasyarakatRepositoryImpl{
		db: db,
	}
}

// DB implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) DB() *gorm.DB {
	return m.db
}

// WithTx implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) WithTx(tx *gorm.DB) IMasyarakatRepository {
	return NewMasyarakatRepositoryImpl(tx)
}

// GetLatestMaterial implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetLatestMaterial(ctx context.Context) ([]*models.Materi, error) {
	var materiList []*models.Materi

	err := m.db.WithContext(ctx).
		Model(&models.Materi{}).
		Preload("MateriGambars").
		Preload("MateriGambars.Gambar").
		Where("status = ?", models.MateriStatusPublished).
		Order("created_at DESC").
		Limit(5).
		Find(&materiList).Error
	if err != nil {
		return nil, err
	}
	return materiList, nil
}

// GetLatestVideo implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetLatestVideo(ctx context.Context) ([]*models.Video, error) {
	var videoList []*models.Video

	err := m.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("status = ?", models.VideoStatusPublished).
		Order("created_at DESC").
		Limit(5).
		Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

// GetAllPublicMateri implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetAllPublicMateri(ctx context.Context, limit int, offset int, search string) ([]*models.Materi, int, error) {
	var (
		materieList []*models.Materi
		count       int64
	)
	if limit <= 0 {
		limit = 10
	}
	query := m.db.WithContext(ctx).
		Model(&models.Materi{}).
		Preload("MateriGambars").
		Preload("MateriGambars.Gambar").
		Where("status = ?", models.MateriStatusPublished)

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&materieList).Error
	if err != nil {
		return nil, 0, err
	}
	return materieList, int(count), nil
}

// GetAllPublicVideo implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetAllPublicVideo(ctx context.Context, limit int, offset int, search string) ([]*models.Video, int, error) {
	var (
		videoList []*models.Video
		count     int64
	)
	if limit <= 0 {
		limit = 10
	}
	query := m.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("status = ?", models.VideoStatusPublished)
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&videoList).Error
	if err != nil {
		return nil, 0, err
	}
	return videoList, int(count), nil
}

// GetPublicMateriByID implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetPublicMateriByID(ctx context.Context, id uuid.UUID) (*models.Materi, error) {
	var materi models.Materi
	err := m.db.WithContext(ctx).
		Model(&models.Materi{}).
		Preload("MateriGambars").
		Preload("MateriGambars.Gambar").
		Where("id = ? AND status = ?", id, models.MateriStatusPublished).
		First(&materi).Error
	if err != nil {
		return nil, err
	}
	return &materi, nil
}

// GetPublicVideoByID implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetPublicVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := m.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ? AND status = ?", id, models.VideoStatusPublished).
		First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// CreateLaporan implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) CreateLaporan(ctx context.Context, laporan *models.Laporan) error {
	laporan.ID = uuid.New()
	return m.db.WithContext(ctx).Create(laporan).Error
}

// CreateGambar implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) CreateGambar(ctx context.Context, gambar *models.Gambar) error {
	gambar.ID = uuid.New()
	return m.db.WithContext(ctx).Create(gambar).Error
}

// CreateLaporanGambar implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) CreateLaporanGambar(ctx context.Context, gambar *models.LaporanGambar) error {
	gambar.ID = uuid.New()
	return m.db.WithContext(ctx).Create(gambar).Error
}

// UpdateLaporan implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) UpdateLaporan(ctx context.Context, laporanId uuid.UUID, updates map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&models.Laporan{}).Where("id = ?", laporanId).Updates(updates).Error
}

func (m *MasyarakatRepositoryImpl) CalculateDFKelurahan(
	ctx context.Context,
	kelurahanId uuid.UUID,
) (float64, string, error) {

	var (
		totalRumah   int64
		rumahPositif int64
		totalWadah   int64
		positifWadah int64
	)

	// 1. Total rumah unik
	err := m.db.WithContext(ctx).
		Table("survey").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kelurahan_id = ?", kelurahanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Distinct("survey.keluarga_id").
		Count(&totalRumah).Error
	if err != nil {
		return 0, "", err
	}

	// 2. Rumah positif (pernah positif)
	err = m.db.WithContext(ctx).
		Table("survey").
		Joins("JOIN survey_item ON survey_item.survey_id = survey.id").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kelurahan_id = ?", kelurahanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Where("survey_item.jumlah_positif > 0").
		Distinct("survey.keluarga_id").
		Count(&rumahPositif).Error
	if err != nil {
		return 0, "", err
	}

	// 3. Total wadah diperiksa
	err = m.db.WithContext(ctx).
		Table("survey_item").
		Select("COALESCE(SUM(jumlah_tempat_air),0)").
		Joins("JOIN survey ON survey.id = survey_item.survey_id").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kelurahan_id = ?", kelurahanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Scan(&totalWadah).Error
	if err != nil {
		return 0, "", err
	}

	// 4. Total wadah positif
	err = m.db.WithContext(ctx).
		Table("survey_item").
		Select("COALESCE(SUM(jumlah_positif),0)").
		Joins("JOIN survey ON survey.id = survey_item.survey_id").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kelurahan_id = ?", kelurahanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Scan(&positifWadah).Error
	if err != nil {
		return 0, "", err
	}

	if totalRumah == 0 || totalWadah == 0 {
		return 0, "Aman", nil
	}

	// Index
	ci := (float64(positifWadah) / float64(totalWadah)) * 100
	hi := (float64(rumahPositif) / float64(totalRumah)) * 100
	bi := (float64(positifWadah) / float64(totalRumah)) * 100

	dfCI := utils.GetDFByCI(ci)
	dfHI := utils.GetDFByHI(hi)
	dfBI := utils.GetDFByBI(bi)

	dfFinal := utils.MaxDF(dfCI, dfHI, dfBI)
	status := utils.GetStatusByDF(dfFinal)

	return float64(dfFinal), status, nil
}

func (m *MasyarakatRepositoryImpl) CalculateDFKecamatan(
	ctx context.Context,
	kecamatanId uuid.UUID,
) (float64, string, error) {

	var (
		totalWadah   int64
		positifWadah int64

		totalRumah   int64
		rumahPositif int64
	)

	// ================================
	// 1. Total Rumah Unik (DISTINCT keluarga)
	// ================================
	err := m.db.WithContext(ctx).
		Table("survey").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kecamatan_id = ?", kecamatanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Distinct("survey.keluarga_id").
		Count(&totalRumah).Error
	if err != nil {
		return 0, "", err
	}

	// ================================
	// 2. Rumah Positif (pernah positif)
	// ================================
	err = m.db.WithContext(ctx).
		Table("survey").
		Joins("JOIN survey_item ON survey_item.survey_id = survey.id").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kecamatan_id = ?", kecamatanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Where("survey_item.jumlah_positif > 0").
		Distinct("survey.keluarga_id").
		Count(&rumahPositif).Error
	if err != nil {
		return 0, "", err
	}

	// ================================
	// 3. Total Wadah diperiksa (SUM)
	// ================================
	err = m.db.WithContext(ctx).
		Table("survey_item").
		Select("COALESCE(SUM(jumlah_tempat_air),0)").
		Joins("JOIN survey ON survey.id = survey_item.survey_id").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kecamatan_id = ?", kecamatanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Scan(&totalWadah).Error
	if err != nil {
		return 0, "", err
	}

	// ================================
	// 4. Total Wadah Positif (SUM)
	// ================================
	err = m.db.WithContext(ctx).
		Table("survey_item").
		Select("COALESCE(SUM(jumlah_positif),0)").
		Joins("JOIN survey ON survey.id = survey_item.survey_id").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Where("keluarga.kecamatan_id = ?", kecamatanId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Scan(&positifWadah).Error
	if err != nil {
		return 0, "", err
	}

	// ================================
	// Validasi kosong
	// ================================
	if totalRumah == 0 || totalWadah == 0 {
		return 0, "Aman", nil
	}

	// ================================
	// Hitung Index Epidemiologi
	// ================================
	ci := (float64(positifWadah) / float64(totalWadah)) * 100
	hi := (float64(rumahPositif) / float64(totalRumah)) * 100
	bi := (float64(positifWadah) / float64(totalRumah)) * 100

	// ================================
	// DF masing-masing index
	// ================================
	dfCI := utils.GetDFByCI(ci)
	dfHI := utils.GetDFByHI(hi)
	dfBI := utils.GetDFByBI(bi)

	// Ambil DF tertinggi
	dfFinal := utils.MaxDF(dfCI, dfHI, dfBI)

	// Status
	status := utils.GetStatusByDF(dfFinal)

	return float64(dfFinal), status, nil
}

// GetLocationResolve implements [IMasyarakatRepository].
func (m *MasyarakatRepositoryImpl) GetLocationResolve(ctx context.Context, kelurahan string, kecamatan string) (*models.CurrenLocation, error) {
	var response models.CurrenLocation

	if kelurahan != "" {
		var dataKelurahan models.Kelurahan
		err := m.db.WithContext(ctx).
			Where("nama_kelurahan ILIKE ?", "%"+kelurahan+"%").
			First(&dataKelurahan).Error

		if err == nil {

			// Hitung DF berdasarkan kelurahan
			df, status, err := m.CalculateDFKelurahan(ctx, dataKelurahan.ID)
			if err != nil {
				return nil, err
			}

			response.Wilayah = dataKelurahan.NamaKelurahan
			response.Scope = "kelurahan"
			response.Df = df
			response.Status = status

			return &response, nil
		}
	}

	if kecamatan != "" {

		var dataKecamatan models.Kecamatan

		err := m.db.WithContext(ctx).
			Where("nama_kecamatan ILIKE ?", "%"+kecamatan+"%").
			First(&dataKecamatan).Error

		if err == nil {

			// Hitung DF berdasarkan kecamatan
			df, status, err := m.CalculateDFKecamatan(ctx, dataKecamatan.ID)
			if err != nil {
				return nil, err
			}

			response.Wilayah = dataKecamatan.NamaKecamatan
			response.Scope = "kecamatan"
			response.Df = df
			response.Status = status

			return &response, nil
		}
	}

	return nil, fmt.Errorf(
		"wilayah tidak ditemukan: kelurahan=%s kecamatan=%s",
		kelurahan,
		kecamatan,
	)
}
