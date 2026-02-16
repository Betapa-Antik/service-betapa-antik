package masyarakatrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMasyarakatRepository interface {
	DB() *gorm.DB
	WithTx(tx *gorm.DB) IMasyarakatRepository
	GetLatestMaterial(ctx context.Context) ([]*models.Materi, error)
	GetLatestVideo(ctx context.Context) ([]*models.Video, error)
	GetAllPublicMateri(ctx context.Context, limit, offset int, search string) ([]*models.Materi, int, error)
	GetPublicMateriByID(ctx context.Context, id uuid.UUID) (*models.Materi, error)
	GetAllPublicVideo(ctx context.Context, limit, offset int, search string) ([]*models.Video, int, error)
	GetPublicVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error)

	CreateLaporan(ctx context.Context, laporan *models.Laporan) error
	CreateGambar(ctx context.Context, gambar *models.Gambar) error
	CreateLaporanGambar(ctx context.Context, gambar *models.LaporanGambar) error
	UpdateLaporan(ctx context.Context, laporanId uuid.UUID, updates map[string]interface{}) error

	GetLocationResolve(ctx context.Context, kelurahan, kecamatan string) (*models.CurrenLocation, error)
}
