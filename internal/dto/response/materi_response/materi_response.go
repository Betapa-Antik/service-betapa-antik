package materiresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type MateriResponse struct {
	ID              uuid.UUID `json:"id"`
	Judul           string    `json:"judul"`
	Deskripsi       string    `json:"deskripsi"`
	Status          string    `json:"status"`
	CatatanTambahan *string   `json:"catatan_tambahan"`
	GambarURLs      []string  `json:"gambar_urls"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

func ToMateriResponse(materi models.Materi) MateriResponse {
	gambars := []string{}
	for _, gambar := range materi.Gambar {
		gambars = append(gambars, gambar.Path)
	}
	return MateriResponse{
		ID:              materi.ID,
		Judul:           materi.Judul,
		Deskripsi:       materi.Deskripsi,
		Status:          materi.Status,
		CatatanTambahan: materi.CatatanTambahan,
		GambarURLs:      gambars,
		CreatedAt:       utils.FormatDate(materi.CreatedAt),
		UpdatedAt:       utils.FormatDate(materi.UpdatedAt),
	}
}
