package roleservice

import (
	rolerequest "betapa-antik-service/internal/dto/request/role_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IRoleService interface {
	Create(ctx context.Context, req *rolerequest.CreateRoleRequest) error
	FindAll(ctx context.Context) ([]models.Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Update(ctx context.Context, id uuid.UUID, req *rolerequest.UpdateRoleRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}
