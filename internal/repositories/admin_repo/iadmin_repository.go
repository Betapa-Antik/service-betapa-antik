package adminrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IAdminRepository interface {
	Register(ctx context.Context, data *models.User) error
	Update(ctx context.Context, id uuid.UUID, data *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	//Manage Petugas
	FindPetugas(ctx context.Context, limit, offset int, search string) ([]*models.User, int, error)
	ApproveOrRejectAkunPetugas(ctx context.Context, petugasId uuid.UUID, status string) error
	ActiveOrNonActiveAkunPetugas(ctx context.Context, petugasId uuid.UUID, status string) error

	GetActiveLupaKataSandi(ctx context.Context, limit, offset int) ([]*models.LupaKataSandi, int, error)
	GetActiveLupaKataSandiById(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error)
	UpdateStatusLupaKataSandi(ctx context.Context, logId uuid.UUID, status string) error

	GetAllLaporan(ctx context.Context, limit, offset int, search string) ([]*models.Laporan, int, error)
	GetLaporanByID(ctx context.Context, laporanId uuid.UUID) (*models.Laporan, error)
	UpdateStatusLaporan(ctx context.Context, laporanId uuid.UUID, updates map[string]interface{}) error

	GetDashboardAdmin(ctx context.Context) (*models.TotalDataDashboardAdmin, error)
	GetSelectKecamatan(ctx context.Context) ([]models.SelectKecamatan, error)
	GetStatistikDFChart(
		ctx context.Context,
		kecamatanId uuid.UUID,
		startDate string,
		endDate string,
	) ([]models.StatistikDFChart, error)
	GetLatestMateri(ctx context.Context) ([]*models.Materi, error)
	GetLatestVideo(ctx context.Context) ([]*models.Video, error)
}
