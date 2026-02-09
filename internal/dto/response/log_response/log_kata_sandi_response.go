package logresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type LogKataSandiResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToLogKataSandiResponse(log models.LupaKataSandi) LogKataSandiResponse {
	return LogKataSandiResponse{
		ID:        log.ID,
		UserID:    log.UserID,
		Status:    log.Status,
		CreatedAt: utils.FormatDate(log.CreatedAt),
		UpdatedAt: utils.FormatDate(log.UpdatedAt),
	}
}

type LupaKataSandiResponse struct {
	ID        uuid.UUID `json:"id"`
	User      string    `json:"user"`
	Puskesmas string    `json:"puskesmas"`
	Role      string    `json:"jabatan"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToLupaKataSandiResponse(lks models.LupaKataSandi) LupaKataSandiResponse {
	return LupaKataSandiResponse{
		ID:        lks.ID,
		User:      lks.User.NamaLengkap,
		Puskesmas: lks.User.Puskesmas.NamaPuskesmas,
		Role:      lks.User.Role.Nama,
		Status:    lks.Status,
		CreatedAt: utils.FormatDate(lks.CreatedAt),
		UpdatedAt: utils.FormatDate(lks.UpdatedAt),
	}
}
