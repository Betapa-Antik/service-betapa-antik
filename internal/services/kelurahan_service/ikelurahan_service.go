package kelurahanservice

import (
	kelurahanrequest "betapa-antik-service/internal/dto/request/kelurahan_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IKelurahanService interface {
	CreateKelurahan(ctx context.Context, req kelurahanrequest.CreateKelurahanRequest) error
	GetAllKelurahan(ctx context.Context, req kelurahanrequest.GetAllKelurahanRequest) ([]*models.Kelurahan, int, error)
	GetKelurahanById(ctx context.Context, kelurahanId uuid.UUID) (*models.Kelurahan, error)
	UpdateKelurahan(ctx context.Context, kelurahanId uuid.UUID, req kelurahanrequest.UpdateKelurahanRequest) error
	DeleteKelurahan(ctx context.Context, kelurahanId uuid.UUID) error
}
