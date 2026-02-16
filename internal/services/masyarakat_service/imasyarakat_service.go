package masyarakatservice

import (
	masyarakatrequest "betapa-antik-service/internal/dto/request/masyarakat_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IMasyarakatService interface {
	GetLatestMaterial(ctx context.Context) ([]*models.Materi, error)
	GetLatestVideo(ctx context.Context) ([]*models.Video, error)
	GetAllPublicMateri(ctx context.Context, req *masyarakatrequest.GetAllPublicMateriRequest) ([]*models.Materi, int, error)
	GetPublicMateriByID(ctx context.Context, materiId uuid.UUID) (*models.Materi, error)
	GetAllPublicVideo(ctx context.Context, req *masyarakatrequest.GetPublicVideoRequest) ([]*models.Video, int, error)
	GetPublicVideoByID(ctx context.Context, videoId uuid.UUID) (*models.Video, error)

	CreateLaporan(ctx context.Context, req masyarakatrequest.CreateLaporanRequest) error
	GetLocationResolve(ctx context.Context, kelurahan, kecamatan string) (*models.CurrenLocation, error)
}
