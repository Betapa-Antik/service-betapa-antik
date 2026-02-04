package kelurahanresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type KelurahanResponse struct {
	ID            uuid.UUID `json:"id"`
	NamaKelurahan string    `json:"nama_kelurahan"`
	KodeKelurahan string    `json:"kode_kelurahan"`
	Kecamatan     *string   `json:"kecamatan,omitempty"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func ToKelurahanResponse(k models.Kelurahan, withKecamatan bool) KelurahanResponse {
	res := KelurahanResponse{
		ID:            k.ID,
		NamaKelurahan: k.NamaKelurahan,
		KodeKelurahan: k.KodeKelurahan,
		CreatedAt:     utils.FormatDate(k.CreatedAt),
		UpdatedAt:     utils.FormatDate(k.UpdatedAt),
	}

	if withKecamatan {
		namaKec := k.Kecamatan.NamaKecamatan
		res.Kecamatan = &namaKec
	}

	return res
}
