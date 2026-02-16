package adminrepo

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepositoryImpl(db *gorm.DB) IAdminRepository {
	return &AdminRepositoryImpl{
		db: db,
	}
}

// Register implements [IAdminRepository].
func (a *AdminRepositoryImpl) Register(ctx context.Context, data *models.User) error {
	data.ID = uuid.New()
	return a.db.WithContext(ctx).Create(data).Error
}

// Update implements [IAdminRepository].
func (a *AdminRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *models.User) error {
	return a.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(data).Error
}

// FindByEmail implements [IAdminRepository].
func (a *AdminRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := a.db.WithContext(ctx).Preload("Role").Joins("LEFT JOIN role ON role.id = \"user\".role_id").Where("email = ? AND status = ? AND role.nama = ?", email, models.UserStatusActive, "ADMIN").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID implements [IAdminRepository].
func (a *AdminRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := a.db.WithContext(ctx).Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ActiveOrNonActiveAkunPetugas implements [IAdminRepository].
func (a *AdminRepositoryImpl) ActiveOrNonActiveAkunPetugas(ctx context.Context, petugasId uuid.UUID, status string) error {
	return a.db.WithContext(ctx).Model(&models.User{}).Where("id = ? AND status in (?, ?)", petugasId, models.UserStatusApproved, models.UserStatusReject).Update("status", status).Error
}

// ApproveOrRejectAkunPetugas implements [IAdminRepository].
func (a *AdminRepositoryImpl) ApproveOrRejectAkunPetugas(ctx context.Context, petugasId uuid.UUID, status string) error {
	return a.db.WithContext(ctx).Model(&models.User{}).Where("id = ? AND status = ?", petugasId, models.UserStatusPending).Update("status", status).Error
}

// FindPetugas implements [IAdminRepository].
func (a *AdminRepositoryImpl) FindPetugas(
	ctx context.Context,
	limit int,
	offset int,
	search string,
) ([]*models.User, int, error) {

	var (
		petugas []*models.User
		total   int64
	)

	if limit <= 0 {
		limit = 10
	}

	searchPattern := "%" + search + "%"

	// ============================
	// 1. COUNT QUERY (FAST)
	// ============================
	countQuery := a.db.WithContext(ctx).
		Model(&models.User{}).
		Joins("LEFT JOIN role ON role.id = \"user\".role_id").
		Joins("LEFT JOIN puskesmas ON puskesmas.id = \"user\".puskesmas_id").
		Where("role.nama = ?", "PETUGAS PUSKESMAS")

	if search != "" {
		countQuery = countQuery.Where(`
			"user".nama_lengkap ILIKE ? OR
			puskesmas.nama_puskesmas ILIKE ?
		`, searchPattern, searchPattern)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ============================
	// 2. DATA QUERY
	// ============================
	dataQuery := a.db.WithContext(ctx).
		Model(&models.User{}).
		Preload("Role").
		Preload("Puskesmas").
		Joins("LEFT JOIN role ON role.id = \"user\".role_id").
		Joins("LEFT JOIN puskesmas ON puskesmas.id = \"user\".puskesmas_id").
		Where("role.nama = ?", "PETUGAS PUSKESMAS")

	if search != "" {
		dataQuery = dataQuery.Where(`
			"user".nama_lengkap ILIKE ? OR
			puskesmas.nama_puskesmas ILIKE ?
		`, searchPattern, searchPattern)
	}

	// Pending di atas
	dataQuery = dataQuery.Order(`
		CASE
			WHEN "user".status = 'pending' THEN 1
			WHEN "user".status = 'approved' THEN 2
			WHEN "user".status = 'active' THEN 3
			WHEN "user".status = 'non-active' THEN 4
			WHEN "user".status = 'reject' THEN 5
			ELSE 6
		END
	`)

	// Terbaru
	dataQuery = dataQuery.Order("\"user\".created_at DESC")

	// Pagination
	if err := dataQuery.
		Limit(limit).
		Offset(offset).
		Find(&petugas).Error; err != nil {
		return nil, 0, err
	}

	return petugas, int(total), nil
}

// GetActiveLupaKataSandi implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetActiveLupaKataSandi(ctx context.Context, limit int, offset int) ([]*models.LupaKataSandi, int, error) {
	var (
		log   []*models.LupaKataSandi
		total int64
	)

	if limit <= 0 {
		limit = 10
	}

	countQuery := a.db.WithContext(ctx).
		Model(&models.LupaKataSandi{})

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := a.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Preload("User.Puskesmas")

	dataQuery = dataQuery.Order(`
		CASE
			WHEN status = 'Pending' THEN 1
			WHEN status = 'Disetujui' THEN 2
			WHEN status = 'Ditolak' THEN 3
			ELSE 4
		END
	`)

	dataQuery = dataQuery.Order("created_at DESC")

	if err := dataQuery.Limit(limit).Offset(offset).Find(&log).Error; err != nil {
		return nil, 0, err
	}

	return log, int(total), nil
}

// GetActiveLupaKataSandiById implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetActiveLupaKataSandiById(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error) {
	var log models.LupaKataSandi
	if err := a.db.WithContext(ctx).Where("id = ?", logId).First(&log).Error; err != nil {
		return nil, err
	}

	return &log, nil
}

// UpdateStatusLupaKataSandi implements [IAdminRepository].
func (a *AdminRepositoryImpl) UpdateStatusLupaKataSandi(ctx context.Context, logId uuid.UUID, status string) error {
	return a.db.WithContext(ctx).Model(&models.LupaKataSandi{}).Where("id = ?", logId).Update("status", status).Error
}

// GetAllLaporan implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetAllLaporan(ctx context.Context, limit int, offset int, search string) ([]*models.Laporan, int, error) {
	var (
		laporanList []*models.Laporan
		count       int64
	)
	if limit <= 0 {
		limit = 10
	}
	query := a.db.WithContext(ctx).Model(&models.Laporan{}).Preload("Puskesmas").Preload("Petugas").Preload("LaporanGambar").Preload("LaporanGambar.Gambar").Where("status  != ?", models.LaporanStatusBaru)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Joins("JOIN puskesmas ON puskesmas.id = laporan.puskesmas_id").
			Joins("JOIN \"user\" ON \"user\".id = laporan.petugas_id").
			Where("puskesmas.nama_puskesmas ILIKE ? OR \"user\".nama_lengkap ILIKE ? OR judul_laporan ILIKE ?", searchPattern, searchPattern, searchPattern)
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&laporanList).Error; err != nil {
		return nil, 0, err
	}
	return laporanList, int(count), nil
}

// GetLaporanByID implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetLaporanByID(ctx context.Context, laporanId uuid.UUID) (*models.Laporan, error) {
	var laporan models.Laporan
	err := a.db.WithContext(ctx).Preload("Puskesmas").Preload("Petugas").Preload("LaporanGambar").Preload("LaporanGambar.Gambar").Where("id = ?", laporanId).First(&laporan).Error
	if err != nil {
		return nil, err
	}
	return &laporan, nil
}

// UpdateStatusLaporan implements [IAdminRepository].
func (a *AdminRepositoryImpl) UpdateStatusLaporan(ctx context.Context, laporanId uuid.UUID, updates map[string]interface{}) error {
	return a.db.WithContext(ctx).Model(&models.Laporan{}).Where("id = ?", laporanId).Updates(updates).Error
}

// GetDashboardAdmin implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetDashboardAdmin(
	ctx context.Context,
) (*models.TotalDataDashboardAdmin, error) {

	var result models.TotalDataDashboardAdmin

	// ================================
	// TOTAL LAPORAN (bukan BARU)
	// ================================
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status != ?", models.LaporanStatusBaru).
		Count(&result.TotalLaporan).Error; err != nil {
		return nil, err
	}

	// ================================
	// HITUNG GROWTH BULAN INI
	// ================================

	// Ambil tanggal awal bulan ini
	now := time.Now()
	startThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Ambil tanggal awal bulan lalu
	startLastMonth := startThisMonth.AddDate(0, -1, 0)

	// Total laporan bulan ini
	var totalThisMonth int64
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status != ?", models.LaporanStatusBaru).
		Where("created_at >= ?", startThisMonth).
		Count(&totalThisMonth).Error; err != nil {
		return nil, err
	}

	// Total laporan bulan lalu
	var totalLastMonth int64
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status != ?", models.LaporanStatusBaru).
		Where("created_at >= ? AND created_at < ?", startLastMonth, startThisMonth).
		Count(&totalLastMonth).Error; err != nil {
		return nil, err
	}

	// ================================
	// HITUNG PERSENTASE GROWTH
	// ================================
	if totalLastMonth == 0 {
		// Kalau bulan lalu belum ada laporan → growth otomatis 100% jika bulan ini ada
		if totalThisMonth > 0 {
			result.GrowthPersenLaporan = 100
		} else {
			result.GrowthPersenLaporan = 0
		}
	} else {
		result.GrowthPersenLaporan =
			(float64(totalThisMonth-totalLastMonth) / float64(totalLastMonth)) * 100
	}

	// ================================
	// Total Kecamatan
	// ================================
	if err := a.db.WithContext(ctx).
		Table("kecamatan").
		Count(&result.TotalKecamatan).Error; err != nil {
		return nil, err
	}

	// Total Puskesmas
	if err := a.db.WithContext(ctx).
		Table("puskesmas").
		Count(&result.TotalPuskesmas).Error; err != nil {
		return nil, err
	}

	// Total Video
	if err := a.db.WithContext(ctx).
		Table("video").
		Count(&result.TotalVideo).Error; err != nil {
		return nil, err
	}

	// Total Materi
	if err := a.db.WithContext(ctx).
		Table("materi").
		Count(&result.TotalMateri).Error; err != nil {
		return nil, err
	}

	// Total Pengajuan
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status = ?", models.LaporanStatusPengajuan).
		Count(&result.TotalLaporanPegajuan).Error; err != nil {
		return nil, err
	}

	// Total Diterima
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status = ?", models.LaporanStatusDisetujui).
		Count(&result.TotalLaporanDiterima).Error; err != nil {
		return nil, err
	}

	// Total Ditolak
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status = ?", models.LaporanStatusDitolak).
		Count(&result.TotalLaporanDitolak).Error; err != nil {
		return nil, err
	}

	// Total Selesai
	if err := a.db.WithContext(ctx).
		Table("laporan").
		Where("status = ?", models.LaporanStatusSelesai).
		Count(&result.TotalLaporanSelesai).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatistikDFChart implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetStatistikDFChart(
	ctx context.Context,
	kecamatanId uuid.UUID,
	startDate string,
	endDate string,
) ([]models.StatistikDFChart, error) {

	var result []models.StatistikDFChart

	query := a.db.WithContext(ctx).
		Table("survey").
		Joins("JOIN keluarga ON keluarga.id = survey.keluarga_id").
		Joins("JOIN kecamatan ON kecamatan.id = keluarga.kecamatan_id").
		Joins("JOIN survey_item ON survey_item.survey_id = survey.id").
		Select(`
    DATE(survey.tanggal) as tanggal,
    kecamatan.nama_kecamatan as kecamatan,

    -- HI
    COALESCE(
        COUNT(DISTINCT CASE 
            WHEN survey_item.jumlah_positif > 0 
            THEN survey.id 
        END)::float 
        / NULLIF(COUNT(DISTINCT survey.id),0) * 100,
    0) as hi,

    -- CI
    COALESCE(
        SUM(survey_item.jumlah_positif)::float
        / NULLIF(SUM(survey_item.jumlah_tempat_air),0) * 100,
    0) as ci,

    -- BI
    COALESCE(
        SUM(survey_item.jumlah_positif)::float
        / NULLIF(COUNT(DISTINCT survey.id),0) * 100,
    0) as bi,

    -- ABJ
    COALESCE(
        (
            COUNT(DISTINCT survey.id)
            -
            COUNT(DISTINCT CASE 
                WHEN survey_item.jumlah_positif > 0 
                THEN survey.id 
            END)
        )::float
        / NULLIF(COUNT(DISTINCT survey.id),0) * 100,
    0) as abj
`)

	// Filter Kecamatan
	if kecamatanId != uuid.Nil {
		query = query.Where("kecamatan.id = ?", kecamatanId)
	}

	// Filter Range Date
	if startDate != "" && endDate != "" {
		query = query.Where("survey.tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	// Grouping
	err := query.
		Where("survey.jenis_survey = ?", models.JenisSurveyJentik).
		Group("DATE(survey.tanggal), kecamatan.nama_kecamatan").
		Order("tanggal ASC").
		Scan(&result).Error
	for i := range result {

		hi := result[i].HI
		ci := result[i].CI
		bi := result[i].BI

		dfHI := utils.GetDFByHI(hi)
		dfCI := utils.GetDFByCI(ci)
		dfBI := utils.GetDFByBI(bi)

		dfFinal := utils.MaxDF(dfHI, dfCI, dfBI)

		result[i].DF = float64(dfFinal)
		result[i].Status = utils.GetStatusByDF(dfFinal)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetSelectKecamatan implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetSelectKecamatan(ctx context.Context) ([]models.SelectKecamatan, error) {
	var result []models.SelectKecamatan

	if err := a.db.WithContext(ctx).
		Table("kecamatan").
		Select(`
			kecamatan.id,
			kecamatan.nama_kecamatan,
			kecamatan.kode_wilayah
		`).Order("kecamatan.nama_kecamatan ASC").
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// GetLatestMateri implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetLatestMateri(ctx context.Context) ([]*models.Materi, error) {
	var materi []*models.Materi
	if err := a.db.WithContext(ctx).Model(&models.Materi{}).
		Preload("MateriGambars").
		Preload("MateriGambars.Gambar").
		Limit(5).
		Order("created_at DESC").
		Find(&materi).Error; err != nil {
		return nil, err
	}

	return materi, nil
}

// GetLatestVideo implements [IAdminRepository].
func (a *AdminRepositoryImpl) GetLatestVideo(ctx context.Context) ([]*models.Video, error) {
	var video []*models.Video
	if err := a.db.WithContext(ctx).Model(&models.Video{}).
		Limit(5).Order("created_at DESC").Find(&video).Error; err != nil {
		return nil, err
	}

	return video, nil
}
