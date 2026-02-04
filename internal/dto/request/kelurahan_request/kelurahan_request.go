package kelurahanrequest

import "github.com/google/uuid"

type CreateKelurahanRequest struct {
	NamaKelurahan string    `json:"nama_kelurahan" form:"nama_kelurahan" validate:"required"`
	KodeKelurahan string    `json:"kode_kelurahan" form:"kode_kelurahan" validate:"required"`
	KecamatanID   uuid.UUID `json:"kecamatan_id" form:"kecamatan_id" validate:"required"`
}

type UpdateKelurahanRequest struct {
	NamaKelurahan string    `json:"nama_kelurahan" form:"nama_kelurahan" validate:"omitempty"`
	KodeKelurahan string    `json:"kode_kelurahan" form:"kode_kelurahan" validate:"omitempty"`
	KecamatanID   uuid.UUID `json:"kecamatan_id" form:"kecamatan_id" validate:"omitempty"`
}

type GetAllKelurahanRequest struct {
	Page        int       `query:"page" validate:"omitempty,min=1"`
	Limit       int       `query:"limit" validate:"omitempty,min=1,max=100"`
	Search      string    `query:"search" validate:"omitempty"`
	KecamatanId uuid.UUID `query:"kecamatanId" validate:"omitempty"`
}
