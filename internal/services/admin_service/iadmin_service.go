package adminservice

import (
	"context"
	"mime/multipart"

	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	"betapa-antik-service/internal/models"

	"github.com/google/uuid"
)

type IAdminService interface {
	Register(ctx context.Context, req *adminrequest.CreateAdminRequest) error
	Login(ctx context.Context, req *authrequest.LoginRequest) (*models.User, string, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, email, nama string) error
	UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, foto *multipart.FileHeader) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Logout(ctx context.Context, token string) error

	//Manage Petugas
	FindPetugas(ctx context.Context, req adminrequest.GetAllPetugasRequest) ([]*models.User, int, error)
	ApproveOrRejectAkunPetugas(ctx context.Context, petugasId uuid.UUID, req adminrequest.UpdateStatusPetugas) error
	ActiveOrNonActiveAkunPetugas(ctx context.Context, petugasId uuid.UUID, req adminrequest.UpdateStatusPetugas) error

	GetActiveLupaKataSandi(ctx context.Context, req adminrequest.GetAllLupaKataSandiRequest) ([]*models.LupaKataSandi, int, error)
	UpdateStatusLupaKataSandi(ctx context.Context, logId uuid.UUID, req adminrequest.UpdateStatusPetugas) error
}
