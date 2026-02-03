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
}
