package laporanresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type LaporanResponse struct {
	ID                uuid.UUID `json:"id"`
	NamaPelapor       string    `json:"nama_pelapor"`
	KontakPelapor     string    `json:"kontak_pelapor"`
	Alamat            string    `json:"alamat"`
	JudulLaporan      string    `json:"judul_laporan"`
	DeskripsiLaporan  string    `json:"deskripsi_laporan"`
	Puskesmas         string    `json:"puskesmas"`
	Petugas           *string   `json:"petugas,omitempty"`
	Status            string    `json:"status"`
	GambarURLs        []any     `json:"gambar_urls"`
	CatatanAdmin      *string   `json:"catatan_admin,omitempty"`
	HasilTindakLanjut *string   `json:"hasil_tindak_lanjut,omitempty"`
	DeliveredAt       *string   `json:"delivered_at,omitempty"`
	DecisionAt        *string   `json:"decision_at,omitempty"`
	FinishedAt        *string   `json:"finished_at,omitempty"`
	CreatedAt         string    `json:"created_at"`
	UpdatedAt         string    `json:"updated_at"`
}

func ToLaporanResponse(laporan models.Laporan) LaporanResponse {
	gambars := []any{}
	for _, gambar := range laporan.LaporanGambar {
		gambars = append(gambars, map[string]interface{}{
			"id":  gambar.ID,          // ID pivot laporan_gambar
			"url": gambar.Gambar.Path, // URL gambar
		})
	}
	return LaporanResponse{
		ID:               laporan.ID,
		NamaPelapor:      laporan.NamaPelapor,
		KontakPelapor:    laporan.KontakPelapor,
		Alamat:           laporan.Alamat,
		JudulLaporan:     laporan.JudulLaporan,
		DeskripsiLaporan: laporan.DeskripsiLaporan,
		Puskesmas:        laporan.Puskesmas.NamaPuskesmas,
		Petugas: func() *string {
			if laporan.PetugasID != uuid.Nil {
				return &laporan.Petugas.NamaLengkap
			} else {
				return nil
			}
		}(),
		Status:     laporan.Status,
		GambarURLs: gambars,
		CatatanAdmin: func() *string {
			if laporan.CatatanAdmin != nil {
				return laporan.CatatanAdmin
			} else {
				return nil
			}
		}(),
		HasilTindakLanjut: func() *string {
			if laporan.HasilTindakLanjut != nil {
				return laporan.HasilTindakLanjut
			} else {
				return nil
			}
		}(),
		DeliveredAt: func() *string {
			if laporan.DeliveredAt != nil {
				str := laporan.DeliveredAt.Format("02-01-2006 15:04")
				return &str
			} else {
				return nil
			}
		}(),
		DecisionAt: func() *string {
			if laporan.DecisionAt != nil {
				str := laporan.DecisionAt.Format("02-01-2006 15:04")
				return &str
			} else {
				return nil
			}
		}(),
		FinishedAt: func() *string {
			if laporan.FinishedAt != nil {
				str := laporan.FinishedAt.Format("02-01-2006 15:04")
				return &str
			} else {
				return nil
			}
		}(),
		CreatedAt: utils.FormatDate(laporan.CreatedAt),
		UpdatedAt: utils.FormatDate(laporan.UpdatedAt),
	}
}
