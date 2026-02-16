package masyarakatrequest

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type GetAllPublicMateriRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}

type GetPublicVideoRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}

type CreateLaporanRequest struct {
	NamaPelapor      string                  `form:"nama_pelapor" validate:"required"`
	KontakPelapor    string                  `form:"kontak_pelapor" validate:"required"`
	Alamat           string                  `form:"alamat" validate:"required"`
	JudulLaporan     string                  `form:"judul_laporan" validate:"required"`
	DeskripsiLaporan string                  `form:"deskripsi_laporan" validate:"required"`
	PuskesmasID      uuid.UUID               `form:"puskesmas_id" validate:"required"`
	Foto             []*multipart.FileHeader `form:"foto,omitempty"`
}
