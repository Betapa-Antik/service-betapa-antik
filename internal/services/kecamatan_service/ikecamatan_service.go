package kecamatanservice

import (
	kecamatanrequest "betapa-antik-service/internal/dto/request/kecamatan_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IKecamatanService interface {
	CreateKecamatan(ctx context.Context, req kecamatanrequest.CreateKecamatanRequest) error
	GetAllKecamatan(ctx context.Context, req kecamatanrequest.GetAllKecamatanRequest) ([]*models.Kecamatan, int, error)
	GetKecamatanById(ctx context.Context, kecamatanId uuid.UUID) (*models.Kecamatan, error)
	UpdateKecamatan(ctx context.Context, kecamatanId uuid.UUID, req kecamatanrequest.UpdateKecamatanRequest) error
	DeleteKecamatan(ctx context.Context, kecamatanId uuid.UUID) error
}
