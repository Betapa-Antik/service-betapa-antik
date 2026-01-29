package roleresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type RoleResponse struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToRoleResponse(role models.Role) RoleResponse {
	return RoleResponse{
		ID:        role.ID,
		Nama:      role.Nama,
		CreatedAt: utils.FormatDate(role.CreatedAt),
		UpdatedAt: utils.FormatDate(role.UpdatedAt),
	}
}
