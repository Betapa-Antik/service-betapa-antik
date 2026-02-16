package adminservice

import (
	"context"
	"mime/multipart"

	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	"betapa-antik-service/internal/models"

	"github.com/google/uuid"
)

type IAdminService interface {
	Register(ctx context.Context, req *adminrequest.CreateAdminRequest) error
	Login(ctx context.Context, req *authrequest.LoginRequest) (*models.User, string, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, email, nama string) error
	UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, foto *multipart.FileHeader) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Logout(ctx context.Context, token string) error

	//Manage Petugas
	FindPetugas(ctx context.Context, req adminrequest.GetAllPetugasRequest) ([]*models.User, int, error)
	ApproveOrRejectAkunPetugas(ctx context.Context, petugasId uuid.UUID, req adminrequest.UpdateStatusPetugas) error
	ActiveOrNonActiveAkunPetugas(ctx context.Context, petugasId uuid.UUID, req adminrequest.UpdateStatusPetugas) error

	GetActiveLupaKataSandi(ctx context.Context, req adminrequest.GetAllLupaKataSandiRequest) ([]*models.LupaKataSandi, int, error)
	UpdateStatusLupaKataSandi(ctx context.Context, logId uuid.UUID, req adminrequest.UpdateStatusPetugas) error

	GetAllLaporan(ctx context.Context, req adminrequest.GetAllLaporanRequest) ([]*models.Laporan, int, error)
	GetLaporanByID(ctx context.Context, laporanId uuid.UUID) (*models.Laporan, error)
	UpdateStatusLaporan(ctx context.Context, laporanId uuid.UUID, req adminrequest.UpdateStatusLaporan) error

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
