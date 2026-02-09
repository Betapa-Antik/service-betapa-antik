package petugasservice

import (
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
)

type IPetugasService interface {
	GetSelectPuskesmas(ctx context.Context, search string) ([]models.SelectPuskesmas, error)
	RegisterAkunPetugas(ctx context.Context, req petugasrequest.RegisterPetugasRequest) error
	LoginPetugas(ctx context.Context, req authrequest.LoginRequest) (*models.User, string, error)
	GetProfilePetugas(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfilePetugas(ctx context.Context, petugasId uuid.UUID, req petugasrequest.UpdatePetugasRequest) error
	LogoutPetugas(ctx context.Context, token string) error

	UbahKataSandi(ctx context.Context, petugasId uuid.UUID, req petugasrequest.UbahKataSandiRequest) error
	LupaKataSandiRequest(ctx context.Context, req petugasrequest.LupaKataSandiRequest) (string, error)
	StatusVerifikasiLupaKataSandi(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error)
	AturUlangKataSandi(ctx context.Context, petugasId uuid.UUID, req petugasrequest.AturUlangKataSandiRequest) error
}
