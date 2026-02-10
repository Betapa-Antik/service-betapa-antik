package keluargaservice

import (
	"betapa-antik-service/configs"
	keluargarequest "betapa-antik-service/internal/dto/request/keluarga_request"
	"betapa-antik-service/internal/models"
	keluargarepo "betapa-antik-service/internal/repositories/keluarga_repo"
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

type KeluargaServiceImpl struct {
	keluargaRepo keluargarepo.IKeluargaRepository
	rdb          *redis.Client
}

func NewKeluargaServiceImpl(keluargaRepo keluargarepo.IKeluargaRepository, rdb *redis.Client) IKeluargaService {
	return &KeluargaServiceImpl{keluargaRepo: keluargaRepo, rdb: rdb}
}

func (k *KeluargaServiceImpl) invalidateKeluargaCache(ctx context.Context, keluargaId uuid.UUID) {
	keyId := fmt.Sprintf("keluarga:%s", keluargaId.String())
	keyAll := "keluarga:all:*"
	iter := k.rdb.Scan(ctx, 0, keyAll, 0).Iterator()
	for iter.Next(ctx) {
		configs.DeleteRedis(ctx, iter.Val())
	}

	_ = configs.DeleteRedis(ctx, keyId)
}

// CreateKeluarga implements [IKeluargaService].
func (k *KeluargaServiceImpl) CreateKeluarga(ctx context.Context, req keluargarequest.CreateKeluargaRequest) error {
	keluarga := &models.Keluarga{
		NamaKepalaKeluarga: req.NamaKepalaKeluarga,
		KecamatanID:        req.KecamatanID,
		KelurahanID:        req.KelurahanID,
		RT:                 req.RT,
		RW:                 req.RW,
		Alamat:             req.Alamat,
	}

	err := utils.RunInTransaction(k.keluargaRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.keluargaRepo.WithTx(tx)

		if err := repoTx.CreateKeluarga(ctx, keluarga); err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat data keluarga", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat data keluarga", 500)
	}

	k.invalidateKeluargaCache(ctx, keluarga.ID)

	return nil
}

// GetAllKeluarga implements [IKeluargaService].
func (k *KeluargaServiceImpl) GetAllKeluarga(ctx context.Context, req keluargarequest.GetAllKeluargaRequest) ([]models.Keluarga, int, error) {
	key := fmt.Sprintf("keluarga:all:search:%s:page:%d:limit:%d", req.Search, req.Page, req.Limit)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var cached cache.CacheKeluarga
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Keluarga, cached.Total, nil
		}
	}

	page := req.Page
	limt := req.Limit
	search := req.Search

	offset := (page - 1) * limt

	data, total, err := k.keluargaRepo.GetAllKeluarga(ctx, limt, offset, search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar keluarga", 500)
	}

	if len(data) == 0 {
		data = []models.Keluarga{}
	}

	cacheData := cache.CacheKeluarga{
		Keluarga: data,
		Total:    total,
	}

	buf, _ := json.Marshal(cacheData)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return data, total, nil
}

// GetKeluargaById implements [IKeluargaService].
func (k *KeluargaServiceImpl) GetKeluargaById(ctx context.Context, keluargaId uuid.UUID) (*models.Keluarga, error) {
	key := fmt.Sprintf("keluarga:%s", keluargaId)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var keluarga models.Keluarga
		if err := json.Unmarshal([]byte(val), &keluarga); err == nil {
			return &keluarga, nil
		}
	}

	keluarga, err := k.keluargaRepo.GetKeluargaById(ctx, keluargaId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil keluarga", 500)
	}

	buf, _ := json.Marshal(keluarga)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return keluarga, nil
}

// UpdateKeluarga implements [IKeluargaService].
func (k *KeluargaServiceImpl) UpdateKeluarga(ctx context.Context, keluargaId uuid.UUID, req keluargarequest.UpdateKeluargaRequest) error {
	return utils.RunInTransaction(k.keluargaRepo.DB(), func(tx *gorm.DB) error {
		repoTx := k.keluargaRepo.WithTx(tx)

		keluarga, err := repoTx.GetKeluargaById(ctx, keluargaId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil keluarga", 500)
		}

		updates := map[string]interface{}{}

		if req.NamaKepalaKeluarga != "" {
			updates["nama_kepala_keluarga"] = req.NamaKepalaKeluarga
		}
		if req.KecamatanID != uuid.Nil {
			updates["kecamatan_id"] = req.KecamatanID
		}
		if req.KelurahanID != uuid.Nil {
			updates["kelurahan_id"] = req.NamaKepalaKeluarga
		}
		if req.RT != "" {
			updates["rt"] = req.RT
		}
		if req.RW != "" {
			updates["rw"] = req.RW
		}
		if req.Alamat != "" {
			updates["alamat"] = req.Alamat
		}

		if len(updates) == 0 {
			return errormessage.NewCustomError(nil, "Tidak ada data untuk diupdate", 400)
		}

		if err := repoTx.UpdateKeluarga(ctx, keluarga.ID, updates); err != nil {
			return errormessage.NewCustomError(err, "Gagal mengupdate keluarga", 500)
		}

		k.invalidateKeluargaCache(ctx, keluarga.ID)
		return nil
	})
}

// DeleteKeluarga implements [IKeluargaService].
func (k *KeluargaServiceImpl) DeleteKeluarga(ctx context.Context, keluargaId uuid.UUID) error {
	keluarga, err := k.keluargaRepo.GetKeluargaById(ctx, keluargaId)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mengambil keluarga", 500)
	}

	if err := k.keluargaRepo.DeleteKeluarga(ctx, keluarga.ID); err != nil {
		return errormessage.NewCustomError(err, "Gagal menghapus keluarga", 500)
	}

	k.invalidateKeluargaCache(ctx, keluarga.ID)
	return nil
}

// GetSelectKecamatan implements [IKeluargaService].
func (k *KeluargaServiceImpl) GetSelectKecamatan(ctx context.Context, search string) ([]models.SelectKecamatan, error) {
	data, err := k.keluargaRepo.GetSelectKecamatan(ctx, search)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil data kecamatan", 500)
	}

	if len(data) == 0 {
		data = []models.SelectKecamatan{}
	}

	return data, nil
}

// GetSelectKelurahan implements [IKeluargaService].
func (k *KeluargaServiceImpl) GetSelectKelurahan(ctx context.Context, kecamatanId uuid.UUID, search string) ([]models.SelectKelurahan, error) {
	data, err := k.keluargaRepo.GetSelectKelurahan(ctx, kecamatanId, search)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil data kelurahan", 500)
	}

	if len(data) == 0 {
		data = []models.SelectKelurahan{}
	}

	return data, nil
}
