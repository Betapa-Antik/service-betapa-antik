package petugasrequest

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type RegisterPetugasRequest struct {
	Foto                *multipart.FileHeader `form:"foto" validate:"required"`
	NamaLengkap         string                `form:"nama_lengkap" validate:"required"`
	NoPegawai           string                `form:"no_pegawai" validate:"required"`
	Email               string                `form:"email" validate:"required,email"`
	PuskesmasId         uuid.UUID             `form:"puskesmas_id" validate:"required"`
	KataSandi           string                `form:"kata_sandi" validate:"required,password"`
	KonfirmasiKataSandi string                `form:"konfirmasi_kata_sandi" validate:"required,eqfield=KataSandi"`
}

type UpdatePetugasRequest struct {
	Foto        *multipart.FileHeader `form:"foto" validate:"omitempty"`
	NamaLengkap string                `form:"nama_lengkap" validate:"omitempty"`
	NoPegawai   string                `form:"no_pegawai" validate:"omitempty"`
	Email       string                `form:"email" validate:"omitempty,email"`
	PuskesmasId uuid.UUID             `form:"puskesmas_id" validate:"omitempty"`
}
