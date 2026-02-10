package keluargaresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type KeluargaResponse struct {
	ID                 uuid.UUID `json:"id"`
	NamaKepalaKeluarga string    `json:"nama_kepala_keluarga"`
	Kecamatan          string    `json:"kecamatan"`
	Kelurahan          string    `json:"kelurahan"`
	RT                 string    `json:"rt"`
	RW                 string    `json:"rw"`
	Alamat             string    `json:"alamat"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

func ToKeluargaResponse(keluarga models.Keluarga) KeluargaResponse {
	return KeluargaResponse{
		ID:                 keluarga.ID,
		NamaKepalaKeluarga: keluarga.NamaKepalaKeluarga,
		Kecamatan:          keluarga.Kecamatan.NamaKecamatan,
		Kelurahan:          keluarga.Kelurahan.NamaKelurahan,
		RT:                 keluarga.RT,
		RW:                 keluarga.RW,
		Alamat:             keluarga.Alamat,
		CreatedAt:          utils.FormatDate(keluarga.CreatedAt),
		UpdatedAt:          utils.FormatDate(keluarga.UpdatedAt),
	}
}
