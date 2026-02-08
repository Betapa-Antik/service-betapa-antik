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
}
