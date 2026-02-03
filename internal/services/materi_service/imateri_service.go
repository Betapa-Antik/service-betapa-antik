package materiservice

import (
	materirequest "betapa-antik-service/internal/dto/request/materi_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IMateriService interface {
	CreateMateri(ctx context.Context, req *materirequest.CreateMateriRequest) error
	GetAllMateri(ctx context.Context, req materirequest.GetAllMateriRequest) ([]*models.Materi, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Materi, error)
	UpdateMateri(ctx context.Context, id uuid.UUID, req *materirequest.UpdateMateriRequest) error
	UpdateStatusMateri(ctx context.Context, id uuid.UUID, req materirequest.UpdateStatusMateriRequest) error
	DeleteMateri(ctx context.Context, id uuid.UUID) error
}
