package petugasrepo

import (
	"betapa-antik-service/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PetugasRepositoryImpl struct {
	db *gorm.DB
}

func NewPetugasRepositoryImpl(db *gorm.DB) IPetugasRepository {
	return &PetugasRepositoryImpl{db: db}
}

// DB implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) DB() *gorm.DB {
	return p.db
}

// WithTx implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) WithTx(tx *gorm.DB) IPetugasRepository {
	return NewPetugasRepositoryImpl(tx)
}

// GetSelectPuskesmas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetSelectPuskesmas(ctx context.Context, search string) ([]models.SelectPuskesmas, error) {
	var result []models.SelectPuskesmas

	query := p.db.WithContext(ctx).
		Table("puskesmas").
		Select(`
			puskesmas.id,
			puskesmas.nama_puskesmas,
			kecamatan.nama_kecamatan,
			kelurahan.nama_kelurahan
		`).
		Joins("LEFT JOIN kecamatan ON kecamatan.id = puskesmas.kecamatan_id").
		Joins("LEFT JOIN kelurahan ON kelurahan.id = puskesmas.kelurahan_id").
		Order("puskesmas.nama_puskesmas ASC")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("puskesmas.nama_puskesmas ILIKE ?", searchPattern)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// GetSelectPuskesmasById implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetSelectPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.SelectPuskesmas, error) {
	var puskesmas models.SelectPuskesmas

	query := p.db.WithContext(ctx).
		Table("puskesmas").
		Select(`
			puskesmas.id,
			puskesmas.nama_puskesmas,
			kecamatan.nama_kecamatan,
			kelurahan.nama_kelurahan
		`).
		Joins("LEFT JOIN kecamatan ON kecamatan.id = puskesmas.kecamatan_id").
		Joins("LEFT JOIN kelurahan ON kelurahan.id = puskesmas.kelurahan_id").
		Order("puskesmas.nama_puskesmas ASC")

	if err := query.Where("puskesmas.id = ?", puskesmasId).First(&puskesmas).Error; err != nil {
		return nil, err
	}

	return &puskesmas, nil
}

// RegisterAkunPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) RegisterAkunPetugas(ctx context.Context, data *models.User) error {
	data.ID = uuid.New()
	return p.db.WithContext(ctx).Create(data).Error
}

// UpdateAkunPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) UpdateAkunPetugas(ctx context.Context, petugasId uuid.UUID, updates map[string]interface{}) error {
	return p.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", petugasId).Updates(updates).Error
}

// FindAkunPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindAkunPetugasById(ctx context.Context, petugasId uuid.UUID) (*models.User, error) {
	var petugas models.User
	if err := p.db.WithContext(ctx).Preload("Role").Preload("Puskesmas").Joins("LEFT JOIN role ON role.id = \"user\".role_id").Where("\"user\".id = ? AND role.nama = ?", petugasId, "PETUGAS PUSKESMAS").First(&petugas).Error; err != nil {
		return nil, err
	}

	return &petugas, nil
}

// FindByEmail implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var petugas models.User
	if err := p.db.WithContext(ctx).Preload("Role").Preload("Puskesmas").Joins("LEFT JOIN role ON role.id = \"user\".role_id").Where("\"user\".email = ? AND role.nama = ?", email, "PETUGAS PUSKESMAS").First(&petugas).Error; err != nil {
		return nil, err
	}

	return &petugas, nil
}

// FindLogForgotPasswordByUserID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindLogForgotPasswordByUserID(ctx context.Context, UserId uuid.UUID) (*models.LupaKataSandi, error) {
	var log models.LupaKataSandi

	err := p.db.WithContext(ctx).
		Where("user_id = ? AND status IN ?", UserId,
			[]string{
				string(models.ForgotPasswordStatusPending),
				string(models.ForgotPasswordStatusPeninjauan),
			},
		).
		First(&log).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &log, nil
}

// CreateLogForgotPassword implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) CreateLogForgotPassword(ctx context.Context, data *models.LupaKataSandi) error {
	data.ID = uuid.New()
	return p.db.WithContext(ctx).Create(data).Error
}

// FindPetugasByEmailPuskesmas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindPetugasByEmailPuskesmas(ctx context.Context, email string, puskesmasId uuid.UUID) (*models.User, error) {
	var user models.User
	if err := p.db.WithContext(ctx).Where("email = ? AND puskesmas_id = ?", email, puskesmasId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindLogForgotPasswordByID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindLogForgotPasswordByID(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error) {
	var log models.LupaKataSandi
	if err := p.db.WithContext(ctx).Where("id = ?", logId).First(&log).Error; err != nil {
		return nil, err
	}

	return &log, nil
}

// UpdateStatusLaporan implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) UpdateStatusLaporan(ctx context.Context, laporanId uuid.UUID, updates map[string]interface{}) error {
	return p.db.WithContext(ctx).Model(&models.Laporan{}).Where("id = ? AND status != ?", laporanId, models.LaporanStatusDitolak).Updates(updates).Error
}

// GetAllLaporan implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetAllLaporan(ctx context.Context, limit int, offset int, search string, petugasId uuid.UUID) ([]*models.Laporan, int, error) {
	var (
		laporanList []*models.Laporan
		count       int64
	)
	if limit <= 0 {
		limit = 10
	}
	query := p.db.WithContext(ctx).Model(&models.Laporan{}).Joins(`JOIN "user" u ON u.puskesmas_id = laporan.puskesmas_id`).
		Where("u.id = ?", petugasId).Preload("Puskesmas").Preload("Petugas").Preload("LaporanGambar").
		Preload("LaporanGambar.Gambar")
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("judul_laporan ILIKE ? OR nama_pelapor ILIKE ?", searchPattern, searchPattern)
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&laporanList).Error; err != nil {
		return nil, 0, err
	}
	return laporanList, int(count), nil
}

// GetLaporanByID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetLaporanByID(ctx context.Context, laporanId uuid.UUID) (*models.Laporan, error) {
	var laporan models.Laporan
	err := p.db.WithContext(ctx).Preload("Puskesmas").Preload("Petugas").Preload("LaporanGambar").Preload("LaporanGambar.Gambar").Where("id = ?", laporanId).First(&laporan).Error
	if err != nil {
		return nil, err
	}
	return &laporan, nil
}

// GetDashboardPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetDashboardPetugas(ctx context.Context, petugasId uuid.UUID) (*models.TotalDataDashboardPetugas, error) {
	var result models.TotalDataDashboardPetugas

	// 1. Hitung Laporan Baru berdasarkan Puskesmas petugas
	err := p.db.WithContext(ctx).
		Model(&models.Laporan{}).
		Joins("JOIN \"user\" ON \"user\".puskesmas_id = laporan.puskesmas_id").
		Where("\"user\".id = ? AND laporan.status = ?", petugasId, models.LaporanStatusBaru).
		Count(&result.TotalLaporanBaru).Error
	if err != nil {
		return nil, err
	}

	// 2. Hitung Total Survey yang dilakukan oleh petugas ini
	err = p.db.WithContext(ctx).
		Model(&models.Survey{}).
		Where("petugas_id = ?", petugasId).
		Count(&result.TotalSurvey).Error
	if err != nil {
		return nil, err
	}

	// 3. Hitung Total Laporan (Semua Status) berdasarkan Puskesmas petugas
	err = p.db.WithContext(ctx).
		Model(&models.Laporan{}).
		Joins("JOIN \"user\" ON \"user\".puskesmas_id = laporan.puskesmas_id").
		Where("\"user\".id = ?", petugasId).
		Count(&result.TotalLaporan).Error
	if err != nil {
		return nil, err
	}

	// 4. Hitung Total Kontainer berdasarkan Puskesmas petugas
	var kontainerList []models.Kontainer
	err = p.db.WithContext(ctx).
		Table("survey_item").
		Select(`
            survey_lokasi.nama_lokasi AS nama_kontainer,
            SUM(survey_item.jumlah_tempat_air) AS jumlah
        `).
		Joins("JOIN survey ON survey.id = survey_item.survey_id").
		Joins("JOIN survey_lokasi ON survey_lokasi.id = survey_item.lokasi_id").
		Joins("JOIN \"user\" ON \"user\".id = survey.petugas_id"). // Join ke user untuk cek puskesmas
		Where("\"user\".puskesmas_id = (SELECT puskesmas_id FROM \"user\" WHERE id = ?)", petugasId).
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Group("survey_lokasi.nama_lokasi").
		Order("jumlah DESC").
		Scan(&kontainerList).Error

	if err != nil {
		return nil, err
	}

	result.TotalKontainer = kontainerList
	return &result, nil
}

// GetLatestLaporanByPuskesmasID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetLatestLaporanByPuskesmasID(ctx context.Context, petugasId uuid.UUID) ([]*models.Laporan, error) {
	var laporanList []*models.Laporan
	err := p.db.WithContext(ctx).
		Model(&models.Laporan{}).
		// Melakukan join ke tabel petugas
		Joins("JOIN \"user\" ON \"user\".puskesmas_id = laporan.puskesmas_id").
		// Filter berdasarkan ID petugas yang sedang login
		Where("\"user\".id = ?", petugasId).
		Order("laporan.created_at DESC").
		Preload("Puskesmas").
		Preload("Petugas").
		Preload("LaporanGambar").
		Preload("LaporanGambar.Gambar").
		Limit(3).
		Find(&laporanList).Error
	if err != nil {
		return nil, err
	}
	return laporanList, nil
}

// GetLatestSurveyByPetugasID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetLatestSurveyByPetugasID(
	ctx context.Context,
	petugasId uuid.UUID,
) ([]*models.Survey, error) {

	var surveys []*models.Survey

	// ================================
	// Subquery: survey terakhir per keluarga
	// ================================
	subQuery := p.db.
		Model(&models.Survey{}).
		Select("DISTINCT ON (keluarga_id) *").
		Where("petugas_id = ?", petugasId).
		Where("jenis_survey = ?", models.JenisSurveyJentik).
		Order("keluarga_id, created_at DESC")

	// ================================
	// Ambil survey + preload keluarga
	// ================================
	err := p.db.WithContext(ctx).
		Table("(?) as survey", subQuery).
		Preload("Keluarga").
		Order("created_at DESC").
		Limit(3).
		Find(&surveys).Error

	if err != nil {
		return nil, err
	}

	// ================================
	// Hitung Status Positif / Negatif
	// ================================
	for _, survey := range surveys {

		var totalPositif int64

		err := p.db.WithContext(ctx).
			Table("survey_item").
			Where("survey_id = ?", survey.ID).
			Select("COALESCE(SUM(jumlah_positif),0)").
			Scan(&totalPositif).Error

		if err != nil {
			return nil, err
		}

		if totalPositif > 0 {
			survey.Status = "Positif"
		} else {
			survey.Status = "Negatif"
		}
	}

	return surveys, nil
}
