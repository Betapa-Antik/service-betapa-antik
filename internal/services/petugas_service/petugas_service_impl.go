package petugasservice

import (
	"betapa-antik-service/configs"
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	"betapa-antik-service/internal/models"
	petugasrepo "betapa-antik-service/internal/repositories/petugas_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	"betapa-antik-service/pkg/workers/producers"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PetugasServiceImpl struct {
	petugasRepo petugasrepo.IPetugasRepository
	roleRepo    rolerepo.IRoleRepository
	rdb         *redis.Client
}

func NewPetugasServiceImpl(petugasRepo petugasrepo.IPetugasRepository, roleRepo rolerepo.IRoleRepository, rdb *redis.Client) IPetugasService {
	return &PetugasServiceImpl{petugasRepo: petugasRepo, roleRepo: roleRepo, rdb: rdb}
}

// GetSelectPuskesmas implements [IPetugasService].
func (p *PetugasServiceImpl) GetSelectPuskesmas(ctx context.Context, search string) ([]models.SelectPuskesmas, error) {
	data, err := p.petugasRepo.GetSelectPuskesmas(ctx, search)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil data puskesmas", 500)
	}

	if len(data) == 0 {
		data = []models.SelectPuskesmas{}
	}

	return data, nil
}

// RegisterAkunPetugas implements [IPetugasService].
func (p *PetugasServiceImpl) RegisterAkunPetugas(ctx context.Context, req petugasrequest.RegisterPetugasRequest) error {
	role, err := p.roleRepo.FindByName(ctx, strings.ToUpper("petugas puskesmas"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "Role admin tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal memeriksa role", 500)
	}

	puskesmas, err := p.petugasRepo.GetSelectPuskesmasById(ctx, req.PuskesmasId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "puskesmas tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal memeriksa puskesmas", 500)
	}

	hashed, err := utils.HashPassword(req.KataSandi)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal generate kata sandi", 500)
	}

	petugas := &models.User{
		Foto:        "",
		NamaLengkap: req.NamaLengkap,
		NoPegawai:   &req.NoPegawai,
		Email:       req.Email,
		PuskesmasID: &puskesmas.ID,
		KataSandi:   hashed,
		RoleID:      role.ID,
		Status:      models.UserStatusPending,
	}

	err = utils.RunInTransaction(p.petugasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.petugasRepo.WithTx(tx)

		if err := repoTx.RegisterAkunPetugas(ctx, petugas); err != nil {
			return errormessage.NewCustomError(err, "Gagal mendaftarkan akun", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mendaftar akun petugas", 500)
	}

	if req.Foto != nil {
		f, err := req.Foto.Open()
		if err == nil {
			ext := filepath.Ext(req.Foto.Filename)
			filename := petugas.ID.String() + ext
			tmpPath := utils.TempFilePath(filename)

			out, err := os.Create(tmpPath)
			if err != nil {
				return nil
			}
			defer out.Close()

			_, _ = io.Copy(out, f)
			pl := payload.PhotoUploadPayload{
				ID:     petugas.ID,
				Folder: "betapa_antik/foto_petugas",
				Files: []payload.PhotoFile{
					{
						Filename: filename,
						Path:     tmpPath,
					},
				},
			}
			producers.PublishPetugasPhotoUpload(pl)
		}
	}
	return nil
}

// LoginPetugas implements [IPetugasService].
func (p *PetugasServiceImpl) LoginPetugas(ctx context.Context, req authrequest.LoginRequest) (*models.User, string, error) {
	allowed, _, ttl, err := utils.IsLoginAllowed(ctx, req.Email)
	if err != nil {
		return nil, "", errormessage.NewCustomError(err, "Gagal memeriksa percobaan login", 500)
	}

	if !allowed {
		return nil, "", errormessage.NewCustomError(errormessage.ErrForbidden, "Terlalu banyak percobaan login, coba lagi setelah "+ttl.String(), 429)
	}

	user, err := p.petugasRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if _, _, incErr := utils.IncrementLoginAttempt(ctx, req.Email); incErr != nil {

		}
		return nil, "", errormessage.NewCustomError(errormessage.ErrUnauthorized, "Email atau kata sandi salah", 401)
	}
	if user.Status != models.UserStatusActive {
		return nil, "", errormessage.NewCustomError(errormessage.ErrForbidden, "Akun Belum diaktifkan", 403)
	}
	if err := utils.CheckPassword(user.KataSandi, req.Password); err != nil {
		count, ttl, incErr := utils.IncrementLoginAttempt(ctx, req.Email)
		if incErr == nil && count >= utils.LoginAttemptLimit {
			return nil, "", errormessage.NewCustomError(errormessage.ErrForbidden, "Terlalu banyak percobaan login, coba lagi setelah"+ttl.String(), 429)
		}
		return nil, "", errormessage.NewCustomError(errormessage.ErrUnauthorized, "Kata Sandi salah", 401)
	}

	if err := utils.ResetLoginAttempt(ctx, req.Email); err != nil {
		// non-fatal, log maybe
	}

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

// GetProfilePetugas implements [IPetugasService].
func (p *PetugasServiceImpl) GetProfilePetugas(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	key := "profile:" + userID.String()
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var u models.User
		if err := json.Unmarshal([]byte(val), &u); err == nil {
			return &u, nil
		}
	}

	petugas, err := p.petugasRepo.FindAkunPetugasById(ctx, userID)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(petugas)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, string(b), 10*time.Minute)
	return petugas, nil
}

// UpdateProfilePetugas implements [IPetugasService].
func (p *PetugasServiceImpl) UpdateProfilePetugas(ctx context.Context, petugasId uuid.UUID, req petugasrequest.UpdatePetugasRequest) error {
	return utils.RunInTransaction(p.petugasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.petugasRepo.WithTx(tx)

		petugas, err := repoTx.FindAkunPetugasById(ctx, petugasId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil petugas", 500)
		}

		updates := map[string]interface{}{}
		if req.NamaLengkap != "" {
			updates["nama_lengkap"] = req.NamaLengkap
		}
		if req.Email != "" {
			updates["email"] = req.Email
		}
		if req.NoPegawai != "" {
			updates["no_pegawai"] = req.NoPegawai
		}
		if req.PuskesmasId != uuid.Nil {
			updates["puskesmas_id"] = req.PuskesmasId
		}

		if err := repoTx.UpdateAkunPetugas(ctx, petugas.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal update profile", 500)
		}

		if req.Foto != nil {
			f, err := req.Foto.Open()
			if err == nil {
				ext := filepath.Ext(req.Foto.Filename)
				filename := petugas.ID.String() + ext
				tmpPath := utils.TempFilePath(filename)

				out, err := os.Create(tmpPath)
				if err != nil {
					return nil
				}
				defer out.Close()

				_, _ = io.Copy(out, f)
				pl := payload.PhotoUploadPayload{
					ID:     petugas.ID,
					Folder: "betapa_antik/foto_petugas",
					Files: []payload.PhotoFile{
						{
							Filename: filename,
							Path:     tmpPath,
						},
					},
				}
				producers.PublishPetugasPhotoUpload(pl)
			}
		}
		_ = configs.DeleteRedis(ctx, "profile:"+petugasId.String())
		return nil
	})
}

// LogoutPetugas implements [IPetugasService].
func (p *PetugasServiceImpl) LogoutPetugas(ctx context.Context, token string) error {
	if token == "" {
		return errormessage.NewCustomError(errormessage.ErrBadRequest, "Token tidak diberikan", 400)
	}
	key := "auth:token:" + token
	if err := configs.DeleteRedis(ctx, key); err != nil {
		return errormessage.NewCustomError(err, "Gagal menghapus token", 500)
	}
	return nil
}

// UbahKataSandi implements [IPetugasService].
func (p *PetugasServiceImpl) UbahKataSandi(ctx context.Context, petugasId uuid.UUID, req petugasrequest.UbahKataSandiRequest) error {
	return utils.RunInTransaction(p.petugasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.petugasRepo.WithTx(tx)

		petugas, err := repoTx.FindAkunPetugasById(ctx, petugasId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil petugas", 500)
		}

		if err := utils.CheckPassword(petugas.KataSandi, req.KataSandiLama); err != nil {
			return errormessage.NewCustomError(errormessage.ErrBadRequest, "Kata sandi lama tidak sesuai", 400)
		}

		hashed, err := utils.HashPassword(req.KataSandiBaru)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal generate kata sandi baru", 500)
		}

		updates := map[string]interface{}{}
		if req.KataSandiBaru != "" {
			updates["kata_sandi"] = hashed
		}

		if err := repoTx.UpdateAkunPetugas(ctx, petugas.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal ubah kata sandi", 500)
		}

		_ = configs.DeleteRedis(ctx, "profile:"+petugasId.String())
		return nil
	})
}

// LupaKataSandiRequest implements [IPetugasService].
func (p *PetugasServiceImpl) LupaKataSandiRequest(ctx context.Context, req petugasrequest.LupaKataSandiRequest) (string, error) {
	petugas, err := p.petugasRepo.FindPetugasByEmailPuskesmas(ctx, req.Email, req.PuskesmasID)
	if err != nil {
		return "", errormessage.NewCustomError(err, "Gagal mengambil data petugas", 500)
	}

	if petugas == nil {
		return "", errormessage.NewCustomError(err, "Akun petugas tidak ditemukan", 500)
	}

	if petugas.Status != models.UserStatusActive {
		return "", errormessage.NewCustomError(err, "Akun anda belum diaktifkan", 500)
	}

	// ✅ CEK REQUEST YANG MASIH AKTIF
	activeReq, err := p.petugasRepo.FindLogForgotPasswordByUserID(ctx, petugas.ID)
	if err != nil {
		return "", errormessage.NewCustomError(err, "Gagal memeriksa permintaan sebelumnya", 500)
	}

	// ✅ Kalau masih ada request aktif → return ID lama (SUCCESS)
	if activeReq != nil {
		return activeReq.ID.String(), nil
	}

	logForgotPassword := &models.LupaKataSandi{
		UserID: petugas.ID,
		Status: models.ForgotPasswordStatusPending,
	}

	err = utils.RunInTransaction(p.petugasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.petugasRepo.WithTx(tx)

		if err := repoTx.CreateLogForgotPassword(ctx, logForgotPassword); err != nil {
			return errormessage.NewCustomError(err, "Gagal melakukan permintaan lupa kata sandi", 500)
		}
		return nil
	})
	if err != nil {
		return "", errormessage.NewCustomError(err, "Gagal melakukan permintaan lupa kata sandi", 500)
	}

	return logForgotPassword.ID.String(), nil
}

// StatusVerifikasiLupaKataSandi implements [IPetugasService].
func (p *PetugasServiceImpl) StatusVerifikasiLupaKataSandi(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error) {
	logLupaKataSandi, err := p.petugasRepo.FindLogForgotPasswordByID(ctx, logId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errormessage.NewCustomError(err, "Permintaan tidak ditemukan", 404)
		}
		return nil, errormessage.NewCustomError(err, "Gagal mengambil permintaan lupa kata sanid", 500)
	}

	return logLupaKataSandi, nil
}

// AturUlangKataSandi implements [IPetugasService].
func (p *PetugasServiceImpl) AturUlangKataSandi(ctx context.Context, petugasId uuid.UUID, req petugasrequest.AturUlangKataSandiRequest) error {
	return utils.RunInTransaction(p.petugasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.petugasRepo.WithTx(tx)

		petugas, err := repoTx.FindAkunPetugasById(ctx, petugasId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil petugas", 500)
		}

		hashed, err := utils.HashPassword(req.KataSandiBaru)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal generate kata sandi baru", 500)
		}

		updates := map[string]interface{}{}
		if req.KataSandiBaru != "" {
			updates["kata_sandi"] = hashed
		}
		if err := repoTx.UpdateAkunPetugas(ctx, petugas.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal ubah kata sandi", 500)
		}

		_ = configs.DeleteRedis(ctx, "profile:"+petugasId.String())
		return nil
	})
}
