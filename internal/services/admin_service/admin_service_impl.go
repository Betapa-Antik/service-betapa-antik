package adminservice

import (
	"context"
	"io"
	"path/filepath"
	"strings"

	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	"betapa-antik-service/internal/models"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	producers "betapa-antik-service/pkg/workers/producers"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminServiceImpl struct {
	adminRepo adminrepo.IAdminRepository
	roleRepo  rolerepo.IRoleRepository
	db        *gorm.DB
}

func NewAdminServiceImpl(adminRepo adminrepo.IAdminRepository, roleRepo rolerepo.IRoleRepository, db *gorm.DB) IAdminService {
	return &AdminServiceImpl{adminRepo: adminRepo, roleRepo: roleRepo, db: db}
}

// Register implements [IAdminService].
func (a *AdminServiceImpl) Register(ctx context.Context, req *adminrequest.CreateAdminRequest) error {
	// find admin role
	role, err := a.roleRepo.FindByName(ctx, strings.ToUpper("admin"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "Role admin tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal memeriksa role", 500)
	}

	// hash password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal generate password", 500)
	}

	user := &models.User{
		Foto:        "",
		NamaLengkap: req.NamaLengkap,
		Email:       req.Email,
		RoleID:      role.ID,
		KataSandi:   hashed,
	}

	// create user inside transaction
	err = utils.RunInTransaction(a.db, func(tx *gorm.DB) error {
		user.ID = uuid.New()
		return tx.WithContext(ctx).Create(user).Error
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat admin", 500)
	}

	// after successful transaction, publish foto upload (async)
	if req.Foto != nil {
		f, err := req.Foto.Open()
		if err == nil {
			data, _ := io.ReadAll(f)
			f.Close()
			// generate filename with uuid + extension
			ext := filepath.Ext(req.Foto.Filename)
			filename := uuid.New().String() + ext
			pl := payload.PhotoUploadPayload{
				UserID: user.ID,
				Folder: "betapa_antik/foto_admin", // specified folder path
				Files: []payload.PhotoFile{
					{
						Filename: filename,
						Data:     data,
					},
				},
			}
			// best-effort publish; do not block on error
			_ = producers.PublishPhotoUpload(pl)
		}
	}

	return nil
}
