package materiservice

import (
	"betapa-antik-service/configs"
	materirequest "betapa-antik-service/internal/dto/request/materi_request"
	"betapa-antik-service/internal/models"
	gambarrepo "betapa-antik-service/internal/repositories/gambar_repo"
	materigambarrepo "betapa-antik-service/internal/repositories/materI_gambar_repo"
	materirepo "betapa-antik-service/internal/repositories/materi_repo"
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

type MateriServiceImpl struct {
	materiRepo   materirepo.IMateriRepository
	gambarRepo   gambarrepo.IGambarRepository
	materiGambar materigambarrepo.IMateriGambarRepository
	rdb          *redis.Client
}

func NewMateriServiceImpl(materiRepo materirepo.IMateriRepository, gambarRepo gambarrepo.IGambarRepository, materiGambar materigambarrepo.IMateriGambarRepository, rdb *redis.Client) IMateriService {
	return &MateriServiceImpl{
		materiRepo:   materiRepo,
		gambarRepo:   gambarRepo,
		materiGambar: materiGambar,
		rdb:          rdb,
	}
}

// InvalidateMateriCache adalah helper untuk menghapus cache terkait materi
func (m *MateriServiceImpl) InvalidateMateriCache(ctx context.Context, id uuid.UUID) {
	// Hapus cache spesifik berdasarkan ID
	_ = configs.DeleteRedis(ctx, "materi:"+id.String())

	// Hapus semua cache list materi (pattern matching)
	iter := m.rdb.Scan(ctx, 0, "materies:all:*", 0).Iterator()
	for iter.Next(ctx) {
		configs.DeleteRedis(ctx, iter.Val())
	}
}

// CreateMateri implements [IMateriService].
func (m *MateriServiceImpl) CreateMateri(ctx context.Context, req *materirequest.CreateMateriRequest) error {
	materi := &models.Materi{
		Judul:           req.Judul,
		Deskripsi:       req.Deskripsi,
		Status:          models.MateriStatusDraft,
		CatatanTambahan: req.CatatanTambahan,
	}

	err := utils.RunInTransaction(m.materiRepo.DB(), func(tx *gorm.DB) error {
		repoTx := m.materiRepo.WithTx(tx)
		if err := repoTx.CreateMateri(ctx, materi); err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat materi", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat materi", 500)
	}

	if len(req.Gambar) > 0 {
		var files []payload.PhotoFile

		for _, v := range req.Gambar {
			src, err := v.Open()
			if err != nil {
				continue
			}
			defer src.Close()

			ext := filepath.Ext(v.Filename)
			filename := uuid.New().String() + ext
			tmpPath := utils.TempFilePath(filename)

			dst, _ := os.Create(tmpPath)
			_, _ = io.Copy(dst, src)
			dst.Close()

			files = append(files, payload.PhotoFile{
				Filename: filename,
				Path:     tmpPath,
			})
		}

		pl := payload.PhotoUploadPayload{
			ID:     materi.ID,
			Folder: "betapa_antik/materi",
			Files:  files,
		}

		// ⬅️ ASYNC, API TIDAK NUNGGU
		producers.PublishMateriPhotoUploadAsync(pl)
	}
	m.InvalidateMateriCache(ctx, materi.ID)

	return nil
}

// GetAllMateri implements [IMateriService].
func (m *MateriServiceImpl) GetAllMateri(ctx context.Context, req materirequest.GetAllMateriRequest) ([]*models.Materi, int, error) {
	key := fmt.Sprintf("materies:all:search:%s:page:%d:limit:%d", req.Search, req.Page, req.Limit)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var cached cache.CacheMateri
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Materies, cached.Total, nil
		}
	}

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

	data, total, err := m.materiRepo.GetAllMateri(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar materi", 500)
	}

	if len(data) == 0 {
		data = []*models.Materi{}
	}

	cacheData := cache.CacheMateri{
		Materies: data,
		Total:    total,
	}

	buf, _ := json.Marshal(cacheData)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)

	return data, total, nil

}

// GetByID implements [IMateriService].
func (m *MateriServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Materi, error) {
	key := fmt.Sprintf("materi:%s", id)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var materi models.Materi
		if err := json.Unmarshal([]byte(val), &materi); err == nil {
			return &materi, nil
		}
	}

	materi, err := m.materiRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
	}

	buf, _ := json.Marshal(materi)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)

	return materi, nil
}

// UpdateMateri implements [IMateriService].
func (m *MateriServiceImpl) UpdateMateri(ctx context.Context, id uuid.UUID, req *materirequest.UpdateMateriRequest) error {
	return utils.RunInTransaction(m.materiRepo.DB(), func(tx *gorm.DB) error {
		repoTx := m.materiRepo.WithTx(tx)
		repoTxGambar := m.gambarRepo.WithTx(tx)
		repoTxMateriGambar := m.materiGambar.WithTx(tx)

		materi, err := repoTx.GetByID(ctx, id)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
		}

		updates := map[string]interface{}{}

		if req.Judul != "" {
			updates["judul"] = req.Judul
		}
		if req.Deskripsi != "" {
			updates["deskripsi"] = req.Deskripsi
		}

		if *req.CatatanTambahan != "" {
			updates["catatan_tambahan"] = req.CatatanTambahan
		}

		if err := repoTx.UpdateMateri(ctx, materi.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal mengupdate materi", 500)
		}
		if len(req.HapusGambarIDs) > 0 {
			var mgIds []uuid.UUID
			for _, v := range req.HapusGambarIDs {
				mgIds = append(mgIds, v)
			}

			data, err := repoTxMateriGambar.FindByIds(ctx, mgIds)
			if err != nil {
				return errormessage.NewCustomError(err, "Gagal mengambil data materi gambar", 500)
			}

			var gambarIDs []uuid.UUID
			var oldPaths []string
			for _, d := range data {
				gambarIDs = append(gambarIDs, d.GambarID)
				oldPaths = append(oldPaths, d.Gambar.Path)
			}

			if err := repoTxMateriGambar.DeleteByIds(ctx, mgIds); err != nil {
				return errormessage.NewCustomError(err, "Gagal menghapus data materi gambar", 500)
			}

			if len(gambarIDs) > 0 {
				if err := repoTxGambar.DeleteByIds(ctx, gambarIDs); err != nil {
					return errormessage.NewCustomError(err, "Gagal menghapus data gambar", 500)
				}
			}

			if len(oldPaths) > 0 {
				producers.PublishDeleteImageAsync(oldPaths)
			}
		}

		if len(req.GambarBaru) > 0 {
			var files []payload.PhotoFile
			for _, g := range req.GambarBaru {
				src, _ := g.Open()
				defer src.Close()

				ext := filepath.Ext(g.Filename)
				filename := uuid.New().String() + ext
				tmpPath := utils.TempFilePath(filename)

				dst, _ := os.Create(tmpPath)
				_, _ = io.Copy(dst, src)
				dst.Close()

				files = append(files, payload.PhotoFile{
					Filename: filename,
					Path:     tmpPath,
				})
			}

			pl := payload.PhotoUploadPayload{
				ID:     id,
				Folder: "betapa_antik/materi",
				Files:  files,
			}

			producers.PublishMateriPhotoUploadAsync(pl)
		}
		m.InvalidateMateriCache(ctx, materi.ID)
		return nil
	})
}

// UpdateStatusMateri implements [IMateriService].
func (m *MateriServiceImpl) UpdateStatusMateri(ctx context.Context, id uuid.UUID, req materirequest.UpdateStatusMateriRequest) error {
	materi, err := m.materiRepo.GetByID(ctx, id)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
	}
	materi.Status = req.Status

	if err := m.materiRepo.UpdateStatusMateri(ctx, materi.ID, req.Status); err != nil {
		return errormessage.NewCustomError(err, "Gagal mengupdate status materi", 500)
	}

	m.InvalidateMateriCache(ctx, materi.ID)
	return nil
}

// DeleteMateri implements [IMateriService].
func (m *MateriServiceImpl) DeleteMateri(ctx context.Context, id uuid.UUID) error {
	return utils.RunInTransaction(m.materiRepo.DB(), func(tx *gorm.DB) error {
		repoMateri := m.materiRepo.WithTx(tx)
		repoGambar := m.gambarRepo.WithTx(tx)
		repoMateriGambar := m.materiGambar.WithTx(tx)

		materi, err := repoMateri.GetByID(ctx, id)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
		}

		// 1️⃣ Ambil relasi materi_gambar
		var mgIds []uuid.UUID
		for _, mg := range materi.MateriGambars {
			mgIds = append(mgIds, mg.ID)
		}

		if len(mgIds) > 0 {
			data, err := repoMateriGambar.FindByIds(ctx, mgIds)
			if err != nil {
				return errormessage.NewCustomError(err, "Gagal mengambil materi gambar", 500)
			}

			var gambarIDs []uuid.UUID
			var oldPaths []string
			for _, d := range data {
				gambarIDs = append(gambarIDs, d.GambarID)
				oldPaths = append(oldPaths, d.Gambar.Path)
			}

			// 2️⃣ Delete pivot table dulu
			if err := repoMateriGambar.DeleteByIds(ctx, mgIds); err != nil {
				return errormessage.NewCustomError(err, "Gagal menghapus materi gambar", 500)
			}

			// 3️⃣ Delete gambar
			if len(gambarIDs) > 0 {
				if err := repoGambar.DeleteByIds(ctx, gambarIDs); err != nil {
					return errormessage.NewCustomError(err, "Gagal menghapus gambar", 500)
				}
			}

			// 4️⃣ Async delete file
			if len(oldPaths) > 0 {
				producers.PublishDeleteImageAsync(oldPaths)
			}
		}

		// 5️⃣ TERAKHIR baru delete materi
		if err := repoMateri.DeleteMateri(ctx, id); err != nil {
			return errormessage.NewCustomError(err, "Gagal menghapus materi", 500)
		}

		m.InvalidateMateriCache(ctx, id)
		return nil
	})
}
