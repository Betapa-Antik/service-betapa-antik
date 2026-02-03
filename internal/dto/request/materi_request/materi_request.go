package materirequest

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type CreateMateriRequest struct {
	Judul           string                  `json:"judul" validate:"required"`
	Deskripsi       string                  `json:"deskripsi" validate:"required"`
	CatatanTambahan *string                 `json:"catatan_tambahan,omitempty"`
	Gambar          []*multipart.FileHeader `json:"gambar,omitempty"`
}

type GetAllMateriRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}

type UpdateMateriRequest struct {
	Judul           string                  `form:"judul" validate:"required"`
	Deskripsi       string                  `form:"deskripsi" validate:"required"`
	CatatanTambahan *string                 `form:"catatan_tambahan,omitempty"`
	GambarBaru      []*multipart.FileHeader `form:"gambar_baru,omitempty"`
	HapusGambarIDs  []uuid.UUID             `form:"hapus_gambar_ids,omitempty"`
}

type UpdateStatusMateriRequest struct {
	Status string `form:"status" validate:"required,oneof=draft published"`
}
