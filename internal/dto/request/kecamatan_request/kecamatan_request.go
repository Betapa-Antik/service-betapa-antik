package kecamatanrequest

import "mime/multipart"

type CreateKecamatanRequest struct {
	Foto          *multipart.FileHeader `form:"foto" validate:"required"`
	NamaKecamatan string                `form:"nama_kecamatan" validate:"required"`
	KodeWilayah   string                `form:"kode_wilayah" validate:"required"`
}
type GetAllKecamatanRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}

type UpdateKecamatanRequest struct {
	Foto          *multipart.FileHeader `form:"foto,omitempty"`
	NamaKecamatan string                `form:"nama_kecamatan" validate:"omitempty"`
	KodeWilayah   string                `form:"kode_wilayah" validate:"omitempty"`
}
