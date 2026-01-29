package rolerepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IRoleRepository interface {
	FindByName(ctx context.Context, nama string) (*models.Role, error)
	Create(ctx context.Context, data *models.Role) error
	FindAll(ctx context.Context) ([]models.Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Update(ctx context.Context, id uuid.UUID, data *models.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}
