package kecamatanresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type KecamatanResponse struct {
	ID            uuid.UUID `json:"id"`
	Foto          string    `json:"foto"`
	NamaKecamatan string    `json:"nama_kecamatan"`
	KodeWilayah   string    `json:"kode_wilayah"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func ToKecamatanResponse(kecamatan models.Kecamatan) KecamatanResponse {
	return KecamatanResponse{
		ID:            kecamatan.ID,
		Foto:          kecamatan.Foto,
		NamaKecamatan: kecamatan.NamaKecamatan,
		KodeWilayah:   kecamatan.KodeWilayah,
		CreatedAt:     utils.FormatDate(kecamatan.CreatedAt),
		UpdatedAt:     utils.FormatDate(kecamatan.UpdatedAt),
	}
}
