package kecamatanservice

import (
	"betapa-antik-service/configs"
	kecamatanrequest "betapa-antik-service/internal/dto/request/kecamatan_request"
	"betapa-antik-service/internal/models"
	kecamatanrepo "betapa-antik-service/internal/repositories/kecamatan_repo"
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

type KecamatanServiceImpl struct {
	kecamatanRepo kecamatanrepo.IKecamatanRepository
	rdb           *redis.Client
}

func NewKecamatanServiceImpl(kecamatanRepo kecamatanrepo.IKecamatanRepository, rdb *redis.Client) IKecamatanService {
	return &KecamatanServiceImpl{kecamatanRepo: kecamatanRepo, rdb: rdb}
}

// InvalidatekecamatanCache adalah helper untuk menghapus cache terkait kecamatan
func (k *KecamatanServiceImpl) InvalidateKecamatanCache(ctx context.Context, id uuid.UUID) {
	// Hapus cache spesifik berdasarkan ID
	_ = configs.DeleteRedis(ctx, "kecamatan:"+id.String())

	// Hapus semua cache list kecamatan (pattern matching)
	iter := k.rdb.Scan(ctx, 0, "kecamatan:all:*", 0).Iterator()
	for iter.Next(ctx) {
		configs.DeleteRedis(ctx, iter.Val())
	}
}

// CreateKecamatan implements [IKecamatanService].
func (k *KecamatanServiceImpl) CreateKecamatan(ctx context.Context, req kecamatanrequest.CreateKecamatanRequest) error {
	kecamatan := &models.Kecamatan{
		Foto:          "",
		NamaKecamatan: req.NamaKecamatan,
		KodeWilayah:   req.KodeWilayah,
	}

	err := utils.RunInTransaction(k.kecamatanRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.kecamatanRepo.WithTx(tx)

		if err := repoTx.CreateKecamatan(ctx, kecamatan); err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat kecamatan", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat kecamatan", 500)
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
				ID:     kecamatan.ID,
				Folder: "betapa_antik/kecamatan",
				Files: []payload.PhotoFile{
					{
						Filename: filename,
						Path:     tmpPath,
					},
				},
			}

			producers.PublishKecamatanPhotoUpload(pl)
		}
	}
	k.InvalidateKecamatanCache(ctx, kecamatan.ID)
	return nil
}

// GetAllKecamatan implements [IKecamatanService].
func (k *KecamatanServiceImpl) GetAllKecamatan(
	ctx context.Context,
	req kecamatanrequest.GetAllKecamatanRequest,
) ([]models.KecamatanWithTotal, int, error) {

	key := fmt.Sprintf("kecamatan:all:search:%s:page:%d:limit:%d",
		req.Search, req.Page, req.Limit,
	)

	// ==========================
	// GET CACHE
	// ==========================
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {

		var cached cache.CacheKecamatan

		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Kecamatans, cached.Total, nil
		}
	}

	// ==========================
	// QUERY DB
	// ==========================
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

	data, total, err := k.kecamatanRepo.GetAllKecamatan(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar kecamatan", 500)
	}

	if len(data) == 0 {
		data = []models.KecamatanWithTotal{}
	}

	// ==========================
	// SET CACHE
	// ==========================
	cacheData := cache.CacheKecamatan{
		Kecamatans: data,
		Total:      total,
	}

	buf, _ := json.Marshal(cacheData)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)

	return data, total, nil
}

// GetKecamatanById implements [IKecamatanService].
func (k *KecamatanServiceImpl) GetKecamatanById(ctx context.Context, kecamatanId uuid.UUID) (*models.KecamatanWithTotal, error) {
	key := fmt.Sprintf("kecamatan:%s", kecamatanId)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var kecamatan models.KecamatanWithTotal
		if err := json.Unmarshal([]byte(val), &kecamatan); err == nil {
			return &kecamatan, nil
		}
	}

	kecamatan, err := k.kecamatanRepo.GetKecamatanById(ctx, kecamatanId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil kecamatan", 500)
	}

	buf, _ := json.Marshal(kecamatan)
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return kecamatan, nil
}

// UpdateKecamatan implements [IKecamatanService].
func (k *KecamatanServiceImpl) UpdateKecamatan(ctx context.Context, kecamatanId uuid.UUID, req kecamatanrequest.UpdateKecamatanRequest) error {
	return utils.RunInTransaction(k.kecamatanRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.kecamatanRepo.WithTx(tx)

		kecamatan, err := repoTx.GetKecamatanById(ctx, kecamatanId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil kecamatan", 500)
		}

		updates := map[string]interface{}{}

		if req.NamaKecamatan != "" {
			updates["nama_kecamatan"] = req.NamaKecamatan
		}
		if req.KodeWilayah != "" {
			updates["kode_wilayah"] = req.KodeWilayah
		}

		if err := repoTx.Update(ctx, kecamatan.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal mengupdate kecamatan", 500)
		}

		if req.Foto != nil {
			var oldPaths []string
			oldPaths = append(oldPaths, kecamatan.Foto)
			if len(oldPaths) > 0 {
				producers.PublishDeleteImageAsync(oldPaths)
			}

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
					ID:     kecamatan.ID,
					Folder: "betapa_antik/kecamatan",
					Files: []payload.PhotoFile{
						{
							Filename: filename,
							Path:     tmpPath,
						},
					},
				}

				producers.PublishKecamatanPhotoUpload(pl)
			}
		}
		k.InvalidateKecamatanCache(ctx, kecamatan.ID)
		return nil
	})
}

// DeleteKecamatan implements [IKecamatanService].
func (k *KecamatanServiceImpl) DeleteKecamatan(ctx context.Context, kecamatanId uuid.UUID) error {
	return utils.RunInTransaction(k.kecamatanRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.kecamatanRepo.WithTx(tx)

		kecamatan, err := repoTx.GetKecamatanById(ctx, kecamatanId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil kecamatan", 500)
		}

		if kecamatan.Foto != "" {
			var oldPaths []string
			oldPaths = append(oldPaths, kecamatan.Foto)
			if len(oldPaths) > 0 {
				producers.PublishDeleteImageAsync(oldPaths)
			}
		}

		if err := repoTx.DeleteKecamatan(ctx, kecamatanId); err != nil {
			return errormessage.NewCustomError(err, "Gagal menghapus kecamatan", 500)
		}

		k.InvalidateKecamatanCache(ctx, kecamatan.ID)
		return nil
	})
}
