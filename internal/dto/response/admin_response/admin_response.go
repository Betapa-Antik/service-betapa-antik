package adminresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type AdminResponse struct {
	ID          uuid.UUID `json:"id"`
	Foto        string    `json:"foto"`
	NamaLengkap string    `json:"nama_lengkap"`
	Email       string    `json:"email"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToAdminResponse(u models.User) AdminResponse {
	return AdminResponse{
		ID:          u.ID,
		Foto:        u.Foto,
		NamaLengkap: u.NamaLengkap,
		Email:       u.Email,
		CreatedAt:   utils.FormatDate(u.CreatedAt),
		UpdatedAt:   utils.FormatDate(u.UpdatedAt),
	}
}
