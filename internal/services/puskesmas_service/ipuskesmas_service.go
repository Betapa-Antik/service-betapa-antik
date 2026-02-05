package puskesmasservice

import (
	puskesmasrequest "betapa-antik-service/internal/dto/request/puskesmas_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IPuskesmasService interface {
	CreatePuskesmas(ctx context.Context, req puskesmasrequest.CreatePuskesmasRequest) error
	GetAllPuskesmas(ctx context.Context, req puskesmasrequest.GetAllKecamatanRequest) ([]models.PuskesmasWithTotal, int, error)
	GetPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.PuskesmasWithTotal, error)
	UpdatePuskesmas(ctx context.Context, puskesmasId uuid.UUID, req puskesmasrequest.UpdatePuskesmasRequest) error
	DeletePuskesmas(ctx context.Context, puskesmasId uuid.UUID) error
}
