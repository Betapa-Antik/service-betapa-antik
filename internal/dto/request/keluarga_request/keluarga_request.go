package keluargarequest

import "github.com/google/uuid"

type CreateKeluargaRequest struct {
	NamaKepalaKeluarga string    `json:"nama_kepala_keluarga" validate:"required"`
	KecamatanID        uuid.UUID `json:"kecamatan_id" validate:"required"`
	KelurahanID        uuid.UUID `json:"kelurahan_id" validate:"required"`
	RT                 string    `json:"rt" validate:"required"`
	RW                 string    `json:"rw" validate:"required"`
	Alamat             string    `json:"alamat" validate:"required"`
}

type UpdateKeluargaRequest struct {
	NamaKepalaKeluarga string    `json:"nama_kepala_keluarga" validate:"omitempty"`
	KecamatanID        uuid.UUID `json:"kecamatan_id" validate:"omitempty"`
	KelurahanID        uuid.UUID `json:"kelurahan_id" validate:"omitempty"`
	RT                 string    `json:"rt" validate:"omitempty"`
	RW                 string    `json:"rw" validate:"omitempty"`
	Alamat             string    `json:"alamat" validate:"omitempty"`
}

type GetAllKeluargaRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}
