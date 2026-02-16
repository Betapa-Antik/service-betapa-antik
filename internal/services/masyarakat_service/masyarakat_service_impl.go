package masyarakatservice

import (
	masyarakatrequest "betapa-antik-service/internal/dto/request/masyarakat_request"
	"betapa-antik-service/internal/models"
	masyarakatrepo "betapa-antik-service/internal/repositories/masyarakat_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"betapa-antik-service/pkg/workers/payload"
	"betapa-antik-service/pkg/workers/producers"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasyarakatServiceImpl struct {
	masyarakatRepo masyarakatrepo.IMasyarakatRepository
}

func NewMasyarakatServiceImpl(masyarakatRepo masyarakatrepo.IMasyarakatRepository) IMasyarakatService {
	return &MasyarakatServiceImpl{
		masyarakatRepo: masyarakatRepo,
	}
}

// GetLatestMaterial implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetLatestMaterial(ctx context.Context) ([]*models.Materi, error) {
	data, err := m.masyarakatRepo.GetLatestMaterial(ctx)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
	}
	if len(data) == 0 {
		data = []*models.Materi{}
	}
	return data, nil
}

// GetLatestVideo implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetLatestVideo(ctx context.Context) ([]*models.Video, error) {
	data, err := m.masyarakatRepo.GetLatestVideo(ctx)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil video", 500)
	}
	if len(data) == 0 {
		data = []*models.Video{}
	}
	return data, nil
}

// GetAllPublicMateri implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetAllPublicMateri(ctx context.Context, req *masyarakatrequest.GetAllPublicMateriRequest) ([]*models.Materi, int, error) {
	data, count, err := m.masyarakatRepo.GetAllPublicMateri(ctx, req.Limit, req.Page, req.Search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
	}
	if len(data) == 0 {
		data = []*models.Materi{}
	}
	return data, count, nil
}

// GetAllPublicVideo implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetAllPublicVideo(ctx context.Context, req *masyarakatrequest.GetPublicVideoRequest) ([]*models.Video, int, error) {
	data, count, err := m.masyarakatRepo.GetAllPublicVideo(ctx, req.Limit, req.Page, req.Search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil video", 500)
	}
	if len(data) == 0 {
		data = []*models.Video{}
	}
	return data, count, nil
}

// GetPublicMateriByID implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetPublicMateriByID(ctx context.Context, materiId uuid.UUID) (*models.Materi, error) {
	data, err := m.masyarakatRepo.GetPublicMateriByID(ctx, materiId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil materi", 500)
	}
	if data == nil {
		return nil, errormessage.NewCustomError(errormessage.ErrNotFound, "Materi tidak ditemukan", 404)
	}
	return data, nil
}

// GetPublicVideoByID implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetPublicVideoByID(ctx context.Context, videoId uuid.UUID) (*models.Video, error) {
	data, err := m.masyarakatRepo.GetPublicVideoByID(ctx, videoId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil video", 500)
	}
	if data == nil {
		return nil, errormessage.NewCustomError(errormessage.ErrNotFound, "Video tidak ditemukan", 404)
	}
	return data, nil
}

// CreateLaporan implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) CreateLaporan(ctx context.Context, req masyarakatrequest.CreateLaporanRequest) error {
	laporan := &models.Laporan{
		NamaPelapor:      req.NamaPelapor,
		KontakPelapor:    req.KontakPelapor,
		Alamat:           req.Alamat,
		JudulLaporan:     req.JudulLaporan,
		DeskripsiLaporan: req.DeskripsiLaporan,
		PuskesmasID:      req.PuskesmasID,
		Status:           models.LaporanStatusBaru,
	}

	err := utils.RunInTransaction(m.masyarakatRepo.DB(), func(tx *gorm.DB) error {
		repoWithTx := m.masyarakatRepo.WithTx(tx)
		err := repoWithTx.CreateLaporan(ctx, laporan)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal membuat laporan", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat laporan", 500)
	}

	if len(req.Foto) > 0 {
		var files []payload.PhotoFile

		for _, v := range req.Foto {
			src, err := v.Open()
			if err != nil {
				continue
			}
			defer src.Close()

			ext := filepath.Ext(v.Filename)
			filename := uuid.New().String() + ext
			tempPath := utils.TempFilePath(filename)

			dst, _ := os.Create(tempPath)
			_, _ = io.Copy(dst, src)
			dst.Close()

			files = append(files, payload.PhotoFile{
				Filename: filename,
				Path:     tempPath,
			})
		}

		pl := payload.PhotoUploadPayload{
			ID:     laporan.ID,
			Folder: "betapa_antik/laporan",
			Files:  files,
		}

		producers.PublishLaporanPhotoUpload(pl)
	}
	return nil
}

// GetLocationResolve implements [IMasyarakatService].
func (m *MasyarakatServiceImpl) GetLocationResolve(ctx context.Context, kelurahan string, kecamatan string) (*models.CurrenLocation, error) {
	if kelurahan == "" && kecamatan == "" {
		return nil, errormessage.NewCustomError(
			errormessage.ErrBadRequest,
			"Kelurahan atau Kecamatan wajib diisi",
			400,
		)
	}

	data, err := m.masyarakatRepo.GetLocationResolve(ctx, kelurahan, kecamatan)
	if err != nil {

		// Kalau error wilayah tidak ditemukan
		if err == gorm.ErrRecordNotFound {
			return nil, errormessage.NewCustomError(
				errormessage.ErrNotFound,
				"Wilayah tidak ditemukan, pastikan data sudah diinput admin",
				404,
			)
		}

		// Error lainnya
		return nil, errormessage.NewCustomError(
			err,
			"Gagal mengambil status DF wilayah",
			500,
		)
	}

	if data.Df == 0 {
		data.Status = "Aman"
	}

	return data, nil
}
