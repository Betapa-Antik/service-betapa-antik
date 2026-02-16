package surveyservice

import (
	"betapa-antik-service/configs"
	surveyrequest "betapa-antik-service/internal/dto/request/survey_request"
	"betapa-antik-service/internal/models"
	surveyrepo "betapa-antik-service/internal/repositories/survey_repo"
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

type SurveyServiceImpl struct {
	surveyRepo surveyrepo.ISurveyRepository
	rdb        *redis.Client
}

func NewSurveyServiceImpl(surveyRepo surveyrepo.ISurveyRepository, rdb *redis.Client) ISurveyService {
	return &SurveyServiceImpl{surveyRepo: surveyRepo, rdb: rdb}
}

func (s *SurveyServiceImpl) InvalidateSurveyCache(ctx context.Context, petugasId uuid.UUID, id uuid.UUID) {
	// Hapus cache spesifik berdasarkan ID
	_ = configs.DeleteRedis(ctx, "survey:"+id.String())

	// Hapus semua cache list materi (pattern matching)
	iter := s.rdb.Scan(ctx, 0, "survey:petugasId:"+petugasId.String()+":all:*", 0).Iterator()
	for iter.Next(ctx) {
		configs.DeleteRedis(ctx, iter.Val())
	}
}

// GetSelectKeluarga implements [ISurveyService].
func (s *SurveyServiceImpl) GetSelectKeluarga(ctx context.Context, search string) ([]models.SelectKeluarga, error) {
	data, err := s.surveyRepo.GetSelectKeluarga(ctx, search)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil data keluarga", 500)
	}

	if len(data) == 0 {
		data = []models.SelectKeluarga{}
	}

	return data, nil
}

// CreateSurvey implements [ISurveyService].
func (s *SurveyServiceImpl) CreateSurvey(ctx context.Context, petugasId uuid.UUID, req surveyrequest.CreateSurveyRequest) error {
	survey := &models.Survey{
		KeluargaID:  req.KeluargaID,
		PetugasID:   petugasId,
		Tanggal:     req.Tanggal,
		JenisSurvey: req.JenisSurvey,
	}

	var itemReqs []surveyrequest.CreateSurveyItemRequest

	if err := utils.DecodeJSON(req.Items, &itemReqs); err != nil {
		return errormessage.NewCustomError(err, "Format items tidak valid", 400)
	}

	err := utils.RunInTransaction(s.surveyRepo.DB(), func(tx *gorm.DB) error {
		repoTx := s.surveyRepo.WithTx(tx)

		if err := repoTx.CreateSurvey(ctx, survey); err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat survey", 500)
		}

		for _, itemReq := range itemReqs {
			lokasi, err := repoTx.FindOrCreateSurveyLokasi(
				ctx,
				itemReq.NamaLokasi,
				req.JenisSurvey,
			)
			if err != nil {
				return err
			}

			item := &models.SurveyItem{
				SurveyID:        survey.ID,
				LokasiID:        lokasi.ID,
				Ditemukan:       itemReq.Ditemukan,
				JumlahTempatAir: itemReq.JumlahTempatAir,
				JumlahNyamuk:    itemReq.JumlahNyamuk,
				JumlahPositif:   itemReq.JumlahPositif,
				JenisPerkiraan:  itemReq.JenisPerkiraan,
				Keterangan:      itemReq.Keterangan,
			}

			if err := repoTx.CreateSurveyItem(ctx, item); err != nil {
				return errormessage.NewCustomError(err, "Gagal membuat survey item", 500)
			}
		}

		if survey.JenisSurvey == models.JenisSurveyJentik {

			if req.FollowUpJentik == nil {
				return errormessage.NewCustomError(
					errormessage.ErrBadRequest,
					"Follow up survey jentik wajib diisi",
					400,
				)
			}

			follow := &models.SurveyFollowUpJentik{
				SurveyID:     survey.ID,
				EdukasiPSN:   req.FollowUpJentik.EdukasiPSN,
				TindakLanjut: req.FollowUpJentik.TindakLanjut,
				Catatan:      req.FollowUpJentik.Catatan,
			}

			if err := repoTx.CreateSurveyFollowUpJentik(ctx, follow); err != nil {
				return errormessage.NewCustomError(err, "Gagal membuat follow up jentik", 500)
			}
		}

		if survey.JenisSurvey == models.JenisSurveyNyamuk {

			if req.FollowUpNyamuk == nil {
				return errormessage.NewCustomError(
					errormessage.ErrBadRequest,
					"Follow up survey nyamuk wajib diisi",
					400,
				)
			}

			follow := &models.SurveyFollowUpNyamuk{
				SurveyID:         survey.ID,
				DitemukanAedes:   req.FollowUpNyamuk.DitemukanAedes,
				TingkatInfestasi: req.FollowUpNyamuk.TingkatInfestasi,
				FoogingStatus:    req.FollowUpNyamuk.FoogingStatus,
				EdukasiOrAbate:   req.FollowUpNyamuk.EdukasiOrAbate,
				Catatan:          req.FollowUpNyamuk.Catatan,
			}

			if err := repoTx.CreateSurveyFollowUpNyamuk(ctx, follow); err != nil {
				return errormessage.NewCustomError(err, "Gagal membuat follow up nyamuk", 500)
			}
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat survey item", 500)
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
			ID:     survey.ID,
			Folder: "betapa_antik/survey",
			Files:  files,
		}

		producers.PublishSurveyPhotoUploadAsync(pl)
	}
	s.InvalidateSurveyCache(ctx, petugasId, survey.ID)
	return nil
}

// GetAllSurvey implements [ISurveyService].
func (s *SurveyServiceImpl) GetAllSurvey(ctx context.Context, req surveyrequest.GetAllSurveyRequest, petugasId uuid.UUID) ([]models.Survey, int, error) {
	key := fmt.Sprintf(
		"survey:petugasId:%s:all:search:%s:page:%d:limit:%d:jenis:%s:start:%s:end:%s",
		petugasId.String(),
		req.Search,
		req.Page,
		req.Limit,
		req.JenisSurvey,
		utils.FormatDateString(req.StartDate),
		utils.FormatDateString(req.EndDate),
	)

	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {

		var cached cache.SurveyCache
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Survey, cached.Total, nil
		}
	}

	page := req.Page
	limit := req.Limit

	offset := (page - 1) * limit

	search := req.Search
	jenisSurvey := req.JenisSurvey
	startDate := req.StartDate
	endDate := req.EndDate

	data, total, err := s.surveyRepo.GetAllSurvey(
		ctx,
		limit,
		offset,
		search,
		jenisSurvey,
		startDate,
		endDate,
		petugasId,
	)

	if err != nil {
		return nil, 0, errormessage.NewCustomError(
			err,
			"Gagal mengambil daftar survey",
			500,
		)
	}

	if len(data) == 0 {
		data = []models.Survey{}
	}

	cacheData := cache.SurveyCache{
		Survey: data,
		Total:  total,
	}

	buf, _ := json.Marshal(cacheData)

	ctxRedis, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_ = configs.SetRedis(ctxRedis, key, buf, time.Minute*10)

	return data, total, nil
}

// GetSurveyByID implements [ISurveyService].
func (s *SurveyServiceImpl) GetSurveyByID(ctx context.Context, surveyId uuid.UUID) (*models.Survey, error) {
	key := fmt.Sprintf("survey:%s", surveyId)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var cached models.Survey
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return &cached, nil
		}
	}

	data, err := s.surveyRepo.GetSurveyByID(ctx, surveyId)
	if err != nil {
		return nil, errormessage.NewCustomError(
			err,
			"Gagal mengambil data survey",
			500,
		)
	}
	buf, _ := json.Marshal(data)
	ctxRedis, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctxRedis, key, buf, time.Minute*10)

	return data, nil
}

// UpdateSurvey implements [ISurveyService].
func (s *SurveyServiceImpl) UpdateSurvey(ctx context.Context, surveyId uuid.UUID, req surveyrequest.UpdateSurveyRequest) error {
	return utils.RunInTransaction(s.surveyRepo.DB(), func(tx *gorm.DB) error {
		repoTx := s.surveyRepo.WithTx(tx)

		survey, err := repoTx.GetSurveyByID(ctx, surveyId)
		if err != nil {
			return errormessage.NewCustomError(
				err,
				"Gagal mengambil data survey",
				500,
			)
		}
		updates := map[string]interface{}{}

		if !req.Tanggal.IsZero() {
			updates["tanggal"] = req.Tanggal
		}
		if req.JenisSurvey != "" {
			updates["jenis_survey"] = req.JenisSurvey
		}
		if req.KeluargaID != uuid.Nil {
			updates["keluarga_id"] = req.KeluargaID
		}
		if len(updates) > 0 {
			if err := repoTx.UpdateSurvey(ctx, survey.ID, updates); err != nil {
				return errormessage.NewCustomError(err, "Gagal mengupdate survey", 500)
			}
		}

		if len(req.HapusItemIDs) > 0 {
			if err := tx.Where("id IN ? AND survey_id = ?", req.HapusItemIDs, surveyId).
				Delete(&models.SurveyItem{}).Error; err != nil {
				return errormessage.NewCustomError(err, "Gagal menghapus item survey", 500)
			}
		}

		if req.Items != "" {
			var itemReqs []surveyrequest.UpdateSurveyItemRequest
			if err := utils.DecodeJSON(req.Items, &itemReqs); err != nil {
				return errormessage.NewCustomError(err, "Format items tidak valid", 400)
			}

			for _, it := range itemReqs {
				itemUpdates := map[string]interface{}{}

				lokasi, err := repoTx.FindOrCreateSurveyLokasi(
					ctx,
					it.NamaLokasi,
					survey.JenisSurvey,
				)
				itemUpdates["lokasi_id"] = lokasi.ID

				itemUpdates["ditemukan"] = it.Ditemukan
				if it.JumlahTempatAir != nil {
					itemUpdates["jumlah_tempat_air"] = *it.JumlahTempatAir
				}
				if it.JumlahPositif != nil {
					itemUpdates["jumlah_positif"] = *it.JumlahPositif
				}
				if it.JumlahNyamuk != nil {
					itemUpdates["jumlah_nyamuk"] = *it.JumlahNyamuk
				}
				if it.JenisPerkiraan != nil {
					itemUpdates["jenis_perkiraan"] = *it.JenisPerkiraan
				}
				if it.Keterangan != nil {
					itemUpdates["keterangan"] = *it.Keterangan
				}

				if err == nil {
					tx.Model(&models.SurveyItem{}).Where("survey_id = ? AND lokasi_id = ?", surveyId, lokasi.ID).Updates(itemUpdates)
				}
			}
		}

		if req.JenisSurvey == models.JenisSurveyJentik && req.FollowUpJentik != nil {
			followUpdates := map[string]interface{}{}
			followUpdates["edukasi_psn"] = req.FollowUpJentik.EdukasiPSN
			if req.FollowUpJentik.TindakLanjut != "" {
				followUpdates["tindak_lanjut"] = req.FollowUpJentik.TindakLanjut
			}
			if req.FollowUpJentik.Catatan != "" {
				followUpdates["catatan"] = req.FollowUpJentik.Catatan
			}

			if err := tx.Model(&models.SurveyFollowUpJentik{}).Where("survey_id = ?", surveyId).Updates(followUpdates).Error; err != nil {
				return errormessage.NewCustomError(err, "Gagal update follow up jentik", 500)
			}
		}

		if req.JenisSurvey == models.JenisSurveyNyamuk && req.FollowUpNyamuk != nil {
			followUpdates := map[string]interface{}{}
			followUpdates["ditemukan_aedes"] = req.FollowUpNyamuk.DitemukanAedes
			followUpdates["edukasi_or_abate"] = req.FollowUpNyamuk.EdukasiOrAbate
			if req.FollowUpNyamuk.TingkatInfestasi != "" {
				followUpdates["tingkat_infestasi"] = req.FollowUpNyamuk.TingkatInfestasi
			}
			if req.FollowUpNyamuk.FoogingStatus != "" {
				followUpdates["fooging_status"] = req.FollowUpNyamuk.FoogingStatus
			}
			if req.FollowUpNyamuk.Catatan != "" {
				followUpdates["catatan"] = req.FollowUpNyamuk.Catatan
			}

			if err := tx.Model(&models.SurveyFollowUpNyamuk{}).Where("survey_id = ?", surveyId).Updates(followUpdates).Error; err != nil {
				return errormessage.NewCustomError(err, "Gagal update follow up nyamuk", 500)
			}
		}
		if len(req.HapusGambarIDs) > 0 {
			var surveyGambar []models.SurveyGambar

			// 1. Ambil data sebelum dihapus untuk mendapatkan Path Cloudinary/Storage
			// Gunakan tx (transaction) yang sedang berjalan
			if err := tx.Preload("Gambar").Where("id IN ?", req.HapusGambarIDs).Find(&surveyGambar).Error; err != nil {
				return errormessage.NewCustomError(err, "Gagal menemukan data gambar", 500)
			}

			if len(surveyGambar) > 0 {
				var gambarIDs []uuid.UUID
				var oldPaths []string

				for _, sg := range surveyGambar {
					gambarIDs = append(gambarIDs, sg.GambarID)
					if sg.Gambar.Path != "" {
						oldPaths = append(oldPaths, sg.Gambar.Path)
					}
				}

				// 2. Hapus record di tabel pivot (SurveyGambar)
				if err := tx.Where("id IN ?", req.HapusGambarIDs).Delete(&models.SurveyGambar{}).Error; err != nil {
					return errormessage.NewCustomError(err, "Gagal menghapus pivot gambar", 500)
				}

				// 3. Hapus record di tabel master (Gambar)
				if len(gambarIDs) > 0 {
					if err := tx.Where("id IN ?", gambarIDs).Delete(&models.Gambar{}).Error; err != nil {
						return errormessage.NewCustomError(err, "Gagal menghapus master gambar", 500)
					}
				}

				// 4. Hapus file fisik via Producer (Cloudinary/S3)
				// Jalankan ini HANYA jika database sukses (setelah return nil atau dipastikan tx commit)
				if len(oldPaths) > 0 {
					producers.PublishDeleteImageAsync(oldPaths)
				}
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
				files = append(files, payload.PhotoFile{Filename: filename, Path: tmpPath})
			}
			producers.PublishSurveyPhotoUploadAsync(payload.PhotoUploadPayload{
				ID: surveyId, Folder: "betapa_antik/survey", Files: files,
			})
		}
		s.InvalidateSurveyCache(ctx, survey.PetugasID, surveyId)
		return nil
	})
}

// DeleteSurvey implements [ISurveyService].
func (s *SurveyServiceImpl) DeleteSurvey(ctx context.Context, surveyId uuid.UUID) error {
	return utils.RunInTransaction(s.surveyRepo.DB(), func(tx *gorm.DB) error {
		repoTx := s.surveyRepo.WithTx(tx)

		// 1. Ambil data survey untuk mendapatkan PetugasID (keperluan cache)
		// dan Preload SurveyGambar.Gambar untuk mendapatkan Path file
		var survey models.Survey
		err := tx.WithContext(ctx).
			Preload("SurveyGambar").
			Preload("SurveyGambar.Gambar").
			Where("id = ?", surveyId).
			First(&survey).Error

		if err != nil {
			return errormessage.NewCustomError(err, "Survey tidak ditemukan", 404)
		}

		// 2. Kumpulkan semua Path gambar dan ID Gambar untuk dihapus
		var oldPaths []string
		var gambarIDs []uuid.UUID
		for _, sg := range survey.SurveyGambar {
			gambarIDs = append(gambarIDs, sg.GambarID)
			if sg.Gambar.Path != "" {
				oldPaths = append(oldPaths, sg.Gambar.Path)
			}
		}

		// 3. Hapus data master gambar
		// (Pivot SurveyGambar akan otomatis terhapus jika Anda menggunakan OnDelete: CASCADE pada database)
		// Jika tidak, hapus manual pivotnya terlebih dahulu.
		if len(gambarIDs) > 0 {
			if err := tx.Where("id IN ?", gambarIDs).Delete(&models.Gambar{}).Error; err != nil {
				return errormessage.NewCustomError(err, "Gagal menghapus data master gambar", 500)
			}
		}

		// 4. Hapus survey utama
		// Dengan GORM OnDelete: CASCADE, SurveyItem dan FollowUp akan ikut terhapus
		if err := repoTx.DeleteSurvey(ctx, surveyId); err != nil {
			return errormessage.NewCustomError(err, "Gagal menghapus data survey", 500)
		}

		// 5. Invalidate Cache

		// 6. Jalankan penghapusan file fisik secara asinkron (Cloudinary/S3)
		if len(oldPaths) > 0 {
			producers.PublishDeleteImageAsync(oldPaths)
		}
		s.InvalidateSurveyCache(ctx, survey.PetugasID, surveyId)
		return nil
	})
}
