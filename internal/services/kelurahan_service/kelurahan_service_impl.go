package kelurahanservice

import (
	"betapa-antik-service/configs"
	kelurahanrequest "betapa-antik-service/internal/dto/request/kelurahan_request"
	"betapa-antik-service/internal/models"
	kelurahanrepo "betapa-antik-service/internal/repositories/kelurahan_repo"
	"betapa-antik-service/pkg/cache"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type KelurahanServiceImpl struct {
	kelurahanRepo kelurahanrepo.IKelurahanRepository
	rdb           *redis.Client
}

func NewKelurahanServiceImpl(kelurahanRepo kelurahanrepo.IKelurahanRepository, rdb *redis.Client) IKelurahanService {
	return &KelurahanServiceImpl{kelurahanRepo: kelurahanRepo, rdb: rdb}
}

func (k *KelurahanServiceImpl) invalidateKelurahanCache(ctx context.Context, kelurahanId uuid.UUID, kecamatanId uuid.UUID) {
	keyId := fmt.Sprintf("kelurahan:%s", kelurahanId.String())
	keyAll := fmt.Sprintf("kelurahan:all:kecamatan:%s:*", kecamatanId.String())
	iter := k.rdb.Scan(ctx, 0, keyAll, 0).Iterator()
	for iter.Next(ctx) {
		configs.DeleteRedis(ctx, iter.Val())
	}

	_ = configs.DeleteRedis(ctx, keyId)

	k.invalidateKecamatanCache(ctx, kecamatanId)
}

func (k *KelurahanServiceImpl) invalidateKecamatanCache(ctx context.Context, kecamatanId uuid.UUID) {
	// hapus cache detail kecamatan
	keyId := fmt.Sprintf("kecamatan:%s", kecamatanId.String())
	_ = configs.DeleteRedis(ctx, keyId)

	// hapus semua cache list kecamatan
	iter := k.rdb.Scan(ctx, 0, "kecamatan:all:*", 0).Iterator()
	for iter.Next(ctx) {
		_ = configs.DeleteRedis(ctx, iter.Val())
	}
}

// CreateKelurahan implements [IKelurahanService].
func (k *KelurahanServiceImpl) CreateKelurahan(ctx context.Context, req kelurahanrequest.CreateKelurahanRequest) error {
	kelurahan := &models.Kelurahan{
		NamaKelurahan: req.NamaKelurahan,
		KodeKelurahan: req.KodeKelurahan,
		KecamatanID:   req.KecamatanID,
	}

	err := utils.RunInTransaction(k.kelurahanRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.kelurahanRepo.WithTx(tx)
		if err := repoTx.CreateKelurahan(ctx, kelurahan); err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat kelurahan", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat kelurahan", 500)
	}

	k.invalidateKelurahanCache(ctx, kelurahan.ID, kelurahan.KecamatanID)

	return nil
}

// GetAllKelurahan implements [IKelurahanService].
func (k *KelurahanServiceImpl) GetAllKelurahan(ctx context.Context, req kelurahanrequest.GetAllKelurahanRequest) ([]models.Kelurahan, int, error) {
	key := fmt.Sprintf("kelurahan:all:kecamatan:%s:search:%s:page:%d:limit:%d", req.KecamatanId, req.Search, req.Page, req.Limit)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var cached cache.CacheKelurahan
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Kelurahans, cached.Total, nil
		}
	}

	page := req.Page
	limit := req.Limit
	search := req.Search
	kecamatanId := req.KecamatanId
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	data, total, err := k.kelurahanRepo.GetAllKelurahan(ctx, kecamatanId, limit, offset, search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar kelurahan", 500)
	}

	if len(data) == 0 {
		data = []models.Kelurahan{}
	}

	cacheData := cache.CacheKelurahan{
		Kelurahans: data,
		Total:      total,
	}

	buf, _ := json.Marshal(cacheData)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return data, total, nil
}

// GetKelurahanById implements [IKelurahanService].
func (k *KelurahanServiceImpl) GetKelurahanById(ctx context.Context, kelurahanId uuid.UUID) (*models.Kelurahan, error) {
	key := fmt.Sprintf("kelurahan:%s", kelurahanId)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var kelurahan models.Kelurahan
		if err := json.Unmarshal([]byte(val), &kelurahan); err == nil {
			return &kelurahan, nil
		}
	}

	kelurahan, err := k.kelurahanRepo.GetKelurahanById(ctx, kelurahanId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil kelurahan", 500)
	}

	buf, _ := json.Marshal(kelurahan)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return kelurahan, nil
}

// UpdateKelurahan implements [IKelurahanService].
func (k *KelurahanServiceImpl) UpdateKelurahan(ctx context.Context, kelurahanId uuid.UUID, req kelurahanrequest.UpdateKelurahanRequest) error {
	return utils.RunInTransaction(k.kelurahanRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.kelurahanRepo.WithTx(tx)

		kelurahan, err := repoTx.GetKelurahanById(ctx, kelurahanId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil kelurahan", 500)
		}

		updates := map[string]interface{}{}

		if req.NamaKelurahan != "" {
			updates["nama_kelurahan"] = req.NamaKelurahan
		}

		if req.KodeKelurahan != "" {
			updates["kode_kelurahan"] = req.KodeKelurahan
		}

		if req.KecamatanID != uuid.Nil {
			updates["kecamatan_id"] = req.KecamatanID
		}

		if len(updates) == 0 {
			return errormessage.NewCustomError(nil, "Tidak ada data untuk diupdate", 400)
		}

		if err := repoTx.UpdateKelurahan(ctx, kelurahan.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal mengupdate kelurahan", 500)
		}
		k.invalidateKelurahanCache(ctx, kelurahan.ID, kelurahan.KecamatanID)
		return nil
	})
}

// DeleteKelurahan implements [IKelurahanService].
func (k *KelurahanServiceImpl) DeleteKelurahan(ctx context.Context, kelurahanId uuid.UUID) error {
	kelurahan, err := k.kelurahanRepo.GetKelurahanById(ctx, kelurahanId)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mengambil kelurahan", 500)
	}

	if err := k.kelurahanRepo.DeleteKelurahan(ctx, kelurahan.ID); err != nil {
		return errormessage.NewCustomError(err, "Gagal menghapus kelurahan", 500)
	}
	k.invalidateKelurahanCache(ctx, kelurahan.ID, kelurahan.KecamatanID)
	return nil
}
