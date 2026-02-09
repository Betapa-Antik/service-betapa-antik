package adminservice

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	"betapa-antik-service/internal/models"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	producers "betapa-antik-service/pkg/workers/producers"

	"errors"

	"betapa-antik-service/configs"

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
		Status:      models.UserStatusActive,
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
			ext := filepath.Ext(req.Foto.Filename)
			filename := uuid.New().String() + ext
			tmpPath := utils.TempFilePath(filename)

			out, err := os.Create(tmpPath)
			if err != nil {
				return nil
			}
			defer out.Close()

			_, _ = io.Copy(out, f)
			pl := payload.PhotoUploadPayload{
				ID:     user.ID,
				Folder: "betapa_antik/foto_admin",
				Files: []payload.PhotoFile{
					{
						Filename: filename,
						Path:     tmpPath,
					},
				},
			}
			// best-effort publish; do not block on error
			producers.PublishPhotoUpload(pl)
		}
	}

	return nil
}

// Login implements [IAdminService].
func (a *AdminServiceImpl) Login(ctx context.Context, req *authrequest.LoginRequest) (*models.User, string, error) {
	// check rate limit
	allowed, _, ttl, err := utils.IsLoginAllowed(ctx, req.Email)
	if err != nil {
		return nil, "", errormessage.NewCustomError(err, "Gagal memeriksa percobaan login", 500)
	}
	if !allowed {
		return nil, "", errormessage.NewCustomError(errormessage.ErrForbidden, "Terlalu banyak percobaan login, coba lagi setelah "+ttl.String(), 429)
	}

	// find user by email
	user, err := a.adminRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// increment attempts on failure
		if _, _, incErr := utils.IncrementLoginAttempt(ctx, req.Email); incErr != nil {
			// log but ignore
		}
		return nil, "", errormessage.NewCustomError(errormessage.ErrUnauthorized, "Email salah", 401)
	}

	// check password
	if err := utils.CheckPassword(user.KataSandi, req.Password); err != nil {
		// increment attempts
		count, ttl, incErr := utils.IncrementLoginAttempt(ctx, req.Email)
		if incErr == nil && count >= utils.LoginAttemptLimit {
			return nil, "", errormessage.NewCustomError(errormessage.ErrForbidden, "Terlalu banyak percobaan login, coba lagi setelah "+ttl.String(), 429)
		}
		return nil, "", errormessage.NewCustomError(errormessage.ErrUnauthorized, "Kata sandi salah", 401)
	}

	// success - reset attempts
	if err := utils.ResetLoginAttempt(ctx, req.Email); err != nil {
		// non-fatal, log maybe
	}

	// generate jwt token
	token, exp, err := utils.GenerateToken(user, 0)
	if err != nil {
		return nil, "", errormessage.NewCustomError(err, "Gagal membuat token", 500)
	}
	// store token in redis with ttl
	ttlDur := time.Until(time.Unix(exp, 0))
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := configs.SetRedis(ctx, "auth:token:"+token, user.ID.String(), ttlDur); err != nil {
		// non-fatal, but return error
		return nil, "", errormessage.NewCustomError(err, "Gagal menyimpan token", 500)
	}

	return user, token, nil
}

// UpdateProfile implements [IAdminService].
func (a *AdminServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, email, nama string) error {
	// fetch user to confirm existence
	user, err := a.adminRepo.FindByID(ctx, userID)
	if err != nil {
		return errormessage.NewCustomError(err, "User tidak ditemukan", 404)
	}

	// if email changed, check uniqueness
	if email != "" && email != user.Email {
		other, err := a.adminRepo.FindByEmail(ctx, email)
		if err == nil {
			if other.ID != userID {
				return errormessage.NewCustomError(errormessage.ErrExists, "Email sudah terdaftar", 400)
			}
		} else {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errormessage.NewCustomError(err, "Gagal memeriksa email", 500)
			}
		}
	}

	data := &models.User{
		Email:       email,
		NamaLengkap: nama,
	}
	if err := a.adminRepo.Update(ctx, userID, data); err != nil {
		return errormessage.NewCustomError(err, "Gagal memperbarui profil", 500)
	}

	// update cache
	updated, err := a.adminRepo.FindByID(ctx, userID)
	if err == nil {
		b, _ := json.Marshal(updated)
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		_ = configs.SetRedis(ctx, "profile:"+userID.String(), string(b), 10*time.Minute)
	}

	return nil
}

// UpdateProfilePhoto implements [IAdminService].
func (a *AdminServiceImpl) UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, foto *multipart.FileHeader) error {
	if foto == nil {
		return errormessage.NewCustomError(errormessage.ErrBadRequest, "File foto tidak ditemukan", 400)
	}
	f, err := foto.Open()
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membaca file foto", 400)
	}
	defer f.Close()
	ext := filepath.Ext(foto.Filename)
	filename := uuid.New().String() + ext
	tmpPath := utils.TempFilePath(filename)

	out, err := os.Create(tmpPath)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal menyimpan file sementara", 500)
	}
	defer out.Close()

	_, _ = io.Copy(out, f)

	pl := payload.PhotoUploadPayload{
		ID:     userID,
		Folder: "betapa_antik/foto_admin",
		Files: []payload.PhotoFile{
			{
				Filename: filename,
				Path:     tmpPath,
			},
		},
	}

	if err := producers.PublishPhotoUpload(pl); err != nil {
		return errormessage.NewCustomError(err, "Gagal memproses unggahan foto", 500)
	}

	// invalidate cache so subsequent GetProfile will fetch fresh data
	_ = configs.DeleteRedis(ctx, "profile:"+userID.String())

	return nil
}

// GetProfile implements [IAdminService]. It tries to read from redis cache first.
func (a *AdminServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	key := "profile:" + userID.String()
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var u models.User
		if err := json.Unmarshal([]byte(val), &u); err == nil {
			return &u, nil
		}
	}

	// fallback to DB
	user, err := a.adminRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(user)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, string(b), 10*time.Minute)
	return user, nil
}

// Logout implements [IAdminService]. It deletes the token key from Redis so it becomes invalid.
func (a *AdminServiceImpl) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errormessage.NewCustomError(errormessage.ErrBadRequest, "Token tidak diberikan", 400)
	}
	key := "auth:token:" + token
	if err := configs.DeleteRedis(ctx, key); err != nil {
		return errormessage.NewCustomError(err, "Gagal menghapus token", 500)
	}
	return nil
}

// ActiveOrNonActiveAkunPetugas implements [IAdminService].
func (a *AdminServiceImpl) ActiveOrNonActiveAkunPetugas(ctx context.Context, petugasId uuid.UUID, req adminrequest.UpdateStatusPetugas) error {
	petugas, err := a.adminRepo.FindByID(ctx, petugasId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "Akun tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal memeriksa Akun", 500)
	}

	petugas.Status = req.Status
	if err := a.adminRepo.ActiveOrNonActiveAkunPetugas(ctx, petugas.ID, req.Status); err != nil {
		return errormessage.NewCustomError(err, "Gagal update status akun petugas", 500)
	}

	_ = configs.DeleteRedis(ctx, "profile:"+petugas.ID.String())

	return nil
}

// ApproveOrRejectAkunPetugas implements [IAdminService].
func (a *AdminServiceImpl) ApproveOrRejectAkunPetugas(ctx context.Context, petugasId uuid.UUID, req adminrequest.UpdateStatusPetugas) error {
	petugas, err := a.adminRepo.FindByID(ctx, petugasId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "Akun tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal memeriksa Akun", 500)
	}

	petugas.Status = req.Status
	if err := a.adminRepo.ApproveOrRejectAkunPetugas(ctx, petugas.ID, req.Status); err != nil {
		return errormessage.NewCustomError(err, "Gagal update status akun petugas", 500)
	}

	_ = configs.DeleteRedis(ctx, "profile:"+petugas.ID.String())

	return nil
}

// FindPetugas implements [IAdminService].
func (a *AdminServiceImpl) FindPetugas(ctx context.Context, req adminrequest.GetAllPetugasRequest) ([]*models.User, int, error) {
	page := req.Page
	limit := req.Limit
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, total, err := a.adminRepo.FindPetugas(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar petugas", 500)
	}

	if len(data) == 0 {
		data = []*models.User{}
	}

	return data, total, nil
}

// GetActiveLupaKataSandi implements [IAdminService].
func (a *AdminServiceImpl) GetActiveLupaKataSandi(ctx context.Context, req adminrequest.GetAllLupaKataSandiRequest) ([]*models.LupaKataSandi, int, error) {
	page := req.Page
	limit := req.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, total, err := a.adminRepo.GetActiveLupaKataSandi(ctx, limit, offset)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar permintaan", 500)
	}

	if len(data) == 0 {
		data = []*models.LupaKataSandi{}
	}

	return data, total, nil
}

// UpdateStatusLupaKataSandi implements [IAdminService].
func (a *AdminServiceImpl) UpdateStatusLupaKataSandi(ctx context.Context, logId uuid.UUID, req adminrequest.UpdateStatusPetugas) error {
	log, err := a.adminRepo.GetActiveLupaKataSandiById(ctx, logId)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
	}
	log.Status = req.Status

	if err := a.adminRepo.UpdateStatusLupaKataSandi(ctx, log.ID, req.Status); err != nil {
		return errormessage.NewCustomError(err, "Gagal mengupdate status", 500)
	}

	return nil
}
