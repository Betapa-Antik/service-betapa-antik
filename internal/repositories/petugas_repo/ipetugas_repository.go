package petugasrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPetugasRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IPetugasRepository
	GetSelectPuskesmas(ctx context.Context, search string) ([]models.SelectPuskesmas, error)
	GetSelectPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.SelectPuskesmas, error)
	RegisterAkunPetugas(ctx context.Context, data *models.User) error
	UpdateAkunPetugas(ctx context.Context, petugasId uuid.UUID, updates map[string]interface{}) error

	FindAkunPetugasById(ctx context.Context, petugasId uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	FindPetugasByEmailPuskesmas(ctx context.Context, email string, puskesmasId uuid.UUID) (*models.User, error)
	CreateLogForgotPassword(ctx context.Context, data *models.LupaKataSandi) error
	FindLogForgotPasswordByUserID(ctx context.Context, UserId uuid.UUID) (*models.LupaKataSandi, error)
	FindLogForgotPasswordByID(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error)

	GetAllLaporan(ctx context.Context, limit, offset int, search string, petugasId uuid.UUID) ([]*models.Laporan, int, error)
	GetLaporanByID(ctx context.Context, laporanId uuid.UUID) (*models.Laporan, error)
	UpdateStatusLaporan(ctx context.Context, laporanId uuid.UUID, updates map[string]interface{}) error

	GetDashboardPetugas(ctx context.Context, petugasId uuid.UUID) (*models.TotalDataDashboardPetugas, error)
	GetLatestLaporanByPuskesmasID(ctx context.Context, petugasId uuid.UUID) ([]*models.Laporan, error)
	GetLatestSurveyByPetugasID(ctx context.Context, petugasId uuid.UUID) ([]*models.Survey, error)
}
