package petugasservice

import (
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	"betapa-antik-service/internal/models"
	petugasrepo "betapa-antik-service/internal/repositories/petugas_repo"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	"betapa-antik-service/pkg/workers/producers"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

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
