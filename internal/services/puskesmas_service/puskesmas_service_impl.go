package puskesmasservice

import (
	"betapa-antik-service/configs"
	puskesmasrequest "betapa-antik-service/internal/dto/request/puskesmas_request"
	"betapa-antik-service/internal/models"
	puskesmasrepo "betapa-antik-service/internal/repositories/puskesmas_repo"
	"betapa-antik-service/pkg/cache"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	"betapa-antik-service/pkg/workers/producers"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PuskesmasServiceImpl struct {
	puskesmasRepo puskesmasrepo.IPuskesmasRepository
	rdb           *redis.Client
}

func NewPuskesmasServiceImpl(puskesmasRepo puskesmasrepo.IPuskesmasRepository, rdb *redis.Client) IPuskesmasService {
	return &PuskesmasServiceImpl{puskesmasRepo: puskesmasRepo, rdb: rdb}
}

func (p *PuskesmasServiceImpl) InvalidatePuskesmasCache(ctx context.Context, id uuid.UUID) {
	// Hapus cache spesifik berdasarkan ID
	_ = configs.DeleteRedis(ctx, "puskesmas:"+id.String())

	// Hapus semua cache list puskesmas (pattern matching)
	iter := p.rdb.Scan(ctx, 0, "puskesmas:all:*", 0).Iterator()
	for iter.Next(ctx) {
		configs.DeleteRedis(ctx, iter.Val())
	}
}

// CreatePuskesmas implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) CreatePuskesmas(ctx context.Context, req puskesmasrequest.CreatePuskesmasRequest) error {
	puskesmas := &models.Puskesmas{
		Foto:          "",
		NamaPuskesmas: req.NamaPuskesmas,
		KecamatanID:   req.KecamatanID,
		KelurahanID:   req.KelurahanID,
		Alamat:        req.Alamat,
		Latitude:      req.Latitude,
		Longtitude:    req.Longtitude,
	}

	err := utils.RunInTransaction(p.puskesmasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.puskesmasRepo.WithTx(tx)

		if err := repoTx.CreatePuskesmas(ctx, puskesmas); err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat puskesmas", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat puskesmas", 500)
	}

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
				ID:     puskesmas.ID,
				Folder: "betapa_antik/puskesmas",
				Files: []payload.PhotoFile{
					{
						Filename: filename,
						Path:     tmpPath,
					},
				},
			}

			producers.PublishPuskesmasPhotoUpload(pl)
		}
	}
	p.InvalidatePuskesmasCache(ctx, puskesmas.ID)
	return nil
}

// GetAllPuskesmas implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) GetAllPuskesmas(ctx context.Context, req puskesmasrequest.GetAllKecamatanRequest) ([]models.PuskesmasWithTotal, int, error) {
	key := fmt.Sprintf("puskesmas:all:search:%s:kecamatan:%s:page:%d:limit:%d", req.Search, req.KecamatanId, req.Page, req.Limit)

	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var cached cache.CachePuskesmas

		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Puskesmas, cached.Total, nil
		}
	}

	page := req.Page
	limit := req.Limit
	search := req.Search
	kecamatan := req.KecamatanId

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, total, err := p.puskesmasRepo.GetAllPuskesmas(ctx, limit, offset, search, kecamatan)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar puskesmas", 500)
	}

	if len(data) == 0 {
		data = []models.PuskesmasWithTotal{}
	}

	cacheData := cache.CachePuskesmas{
		Puskesmas: data,
		Total:     total,
	}

	buf, _ := json.Marshal(cacheData)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)

	return data, total, nil
}

// GetPuskesmasById implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) GetPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.PuskesmasWithTotal, error) {
	key := fmt.Sprintf("puskesmas:%s", puskesmasId)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var puskesmas models.PuskesmasWithTotal
		if err := json.Unmarshal([]byte(val), &puskesmas); err == nil {
			return &puskesmas, nil
		}
	}

	puskesmas, err := p.puskesmasRepo.GetPuskesmasById(ctx, puskesmasId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil puskesmas", 500)
	}

	buf, _ := json.Marshal(puskesmas)
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return puskesmas, nil
}

// UpdatePuskesmas implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) UpdatePuskesmas(ctx context.Context, puskesmasId uuid.UUID, req puskesmasrequest.UpdatePuskesmasRequest) error {
	return utils.RunInTransaction(p.puskesmasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.puskesmasRepo.WithTx(tx)

		puskesmas, err := repoTx.GetPuskesmasById(ctx, puskesmasId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil puskesmas", 500)
		}

		updates := map[string]interface{}{}

		if req.NamaPuskesmas != "" {
			updates["nama_puskesmas"] = req.NamaPuskesmas
		}
		if req.KecamatanID != uuid.Nil {
			updates["kecamatan_id"] = req.KecamatanID
		}
		if req.KelurahanID != uuid.Nil {
			updates["kelurahan_id"] = req.KelurahanID
		}
		if req.Alamat != "" {
			updates["alamat"] = req.Alamat
		}
		if req.Latitude != "" {
			updates["latitude"] = req.Latitude
		}
		if req.Longtitude != "" {
			updates["longitude"] = req.Longtitude
		}

		if err := repoTx.UpdatePuskesmas(ctx, puskesmas.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal mengupdate puskesmas", 500)
		}

		if req.Foto != nil {
			var oldPaths []string
			oldPaths = append(oldPaths, puskesmas.Foto)

			if len(oldPaths) > 0 {
				producers.PublishDeleteImageAsync(oldPaths)
			}

			f, err := req.Foto.Open()
			if err == nil {
				ext := filepath.Ext(req.Foto.Filename)
				filename := puskesmas.ID.String() + ext
				tmpPath := utils.TempFilePath(filename)
				out, err := os.Create(tmpPath)
				if err != nil {
					return nil
				}
				defer out.Close()

				_, _ = io.Copy(out, f)
				pl := payload.PhotoUploadPayload{
					ID:     puskesmas.ID,
					Folder: "betapa_antik/puskesmas",
					Files: []payload.PhotoFile{
						{
							Filename: filename,
							Path:     tmpPath,
						},
					},
				}

				producers.PublishPuskesmasPhotoUpload(pl)
			}
		}
		p.InvalidatePuskesmasCache(ctx, puskesmas.ID)
		return nil
	})
}

// DeletePuskesmas implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) DeletePuskesmas(ctx context.Context, puskesmasId uuid.UUID) error {
	return utils.RunInTransaction(p.puskesmasRepo.DB(), func(tx *gorm.DB) error {
		repoTx := p.puskesmasRepo.WithTx(tx)

		puskesmas, err := repoTx.GetPuskesmasById(ctx, puskesmasId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal menghapus puskesmas", 500)
		}

		if puskesmas.Foto != "" {
			var oldPaths []string
			oldPaths = append(oldPaths, puskesmas.Foto)

			if len(oldPaths) > 0 {
				producers.PublishDeleteImageAsync(oldPaths)
			}
		}

		if err := repoTx.DeletePuskesmas(ctx, puskesmasId); err != nil {
			return errormessage.NewCustomError(err, "Gagal menghapus puskesmas", 500)
		}

		p.InvalidatePuskesmasCache(ctx, puskesmas.ID)
		return nil
	})
}

// GetSelectKecamatan implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error) {
	data, err := p.puskesmasRepo.GetSelectKecamatan(ctx, search)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil data kecamatan", 500)
	}

	if len(data) == 0 {
		data = []models.SelectKecamatan{}
	}

	return data, nil
}

// GetSelectKelurahan implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error) {
	data, err := p.puskesmasRepo.GetSelectKelurahan(ctx, kecamatanId, search)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil data kelurahan", 500)
	}

	if len(data) == 0 {
		data = []models.SelectKelurahan{}
	}

	return data, nil
}

// GetPetugasByPuskesmasId implements [IPuskesmasService].
func (p *PuskesmasServiceImpl) GetPetugasByPuskesmasId(ctx context.Context, puskesmasId uuid.UUID, search string) ([]*models.PetugasWithTotalSurvey, error) {
	return p.puskesmasRepo.GetAllPetugasByPuskesmasId(ctx, puskesmasId, search)
}
