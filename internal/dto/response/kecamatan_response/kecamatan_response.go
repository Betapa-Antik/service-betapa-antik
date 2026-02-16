package kecamatanresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type KecamatanResponse struct {
	ID             uuid.UUID `json:"id"`
	Foto           string    `json:"foto"`
	NamaKecamatan  string    `json:"nama_kecamatan"`
	KodeWilayah    string    `json:"kode_wilayah"`
	TotalKelurahan int       `json:"total_kelurahan"`
	TotalPuskesmas int       `json:"total_puskesmas"`
	HI             float64   `json:"hi"`
	CI             float64   `json:"ci"`
	BI             float64   `json:"bi"`
	ABJ            float64   `json:"abj"`
	DF             int       `json:"df"`
	Status         string    `json:"status"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func ToKecamatanResponse(kecamatan models.KecamatanWithTotal) KecamatanResponse {
	return KecamatanResponse{
		ID:             kecamatan.ID,
		Foto:           kecamatan.Foto,
		NamaKecamatan:  kecamatan.NamaKecamatan,
		KodeWilayah:    kecamatan.KodeWilayah,
		TotalKelurahan: kecamatan.TotalKelurahan,
		TotalPuskesmas: kecamatan.TotalPuskesmas,
		HI:             kecamatan.HI,
		CI:             kecamatan.CI,
		BI:             kecamatan.BI,
		ABJ:            kecamatan.ABJ,
		DF:             kecamatan.DF,
		Status:         kecamatan.Status,
		CreatedAt:      utils.FormatDate(kecamatan.CreatedAt),
		UpdatedAt:      utils.FormatDate(kecamatan.UpdatedAt),
	}
}

type KecamatanSelectedResponse struct {
	ID            uuid.UUID `json:"id"`
	NamaKecamatan string    `json:"nama_kecamatan"`
	KodeWilayah   string    `json:"kode_wilayah"`
}

func ToKecamatanSelectedResponse(kecamatan models.SelectKecamatan) KecamatanSelectedResponse {
	return KecamatanSelectedResponse{
		ID:            kecamatan.ID,
		NamaKecamatan: kecamatan.NamaKecamatan,
		KodeWilayah:   kecamatan.KodeWilayah,
	}
}
