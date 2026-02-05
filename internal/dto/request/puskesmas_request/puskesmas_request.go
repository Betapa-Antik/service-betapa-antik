package puskesmasrequest

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type CreatePuskesmasRequest struct {
	Foto          *multipart.FileHeader `form:"foto" validate:"required"`
	NamaPuskesmas string                `form:"nama_puskesmas" validate:"required"`
	KecamatanID   uuid.UUID             `form:"kecamatan_id" validate:"required"`
	KelurahanID   uuid.UUID             `form:"kelurahan_id" validate:"required"`
	Alamat        string                `form:"alamat" validate:"required"`
	Latitude      string                `form:"latitude" validate:"required"`
	Longtitude    string                `form:"longitude" validate:"required"`
}

type GetAllKecamatanRequest struct {
	Page        int       `query:"page" validate:"omitempty,min=1"`
	Limit       int       `query:"limit" validate:"omitempty,min=1,max=100"`
	Search      string    `query:"search" validate:"omitempty"`
	KecamatanId uuid.UUID `query:"kecamatan_id" validate:"omitempty"`
}

type UpdatePuskesmasRequest struct {
	Foto          *multipart.FileHeader `form:"foto" validate:"omitempty"`
	NamaPuskesmas string                `form:"nama_puskesmas" validate:"omitempty"`
	KecamatanID   uuid.UUID             `form:"kecamatan_id" validate:"omitempty"`
	KelurahanID   uuid.UUID             `form:"kelurahan_id" validate:"omitempty"`
	Alamat        string                `form:"alamat" validate:"omitempty"`
	Latitude      string                `form:"latitude" validate:"omitempty"`
	Longtitude    string                `form:"longitude" validate:"omitempty"`
}
