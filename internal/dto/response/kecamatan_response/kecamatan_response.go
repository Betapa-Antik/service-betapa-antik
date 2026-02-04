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
		CreatedAt:      utils.FormatDate(kecamatan.CreatedAt),
		UpdatedAt:      utils.FormatDate(kecamatan.UpdatedAt),
	}
}
