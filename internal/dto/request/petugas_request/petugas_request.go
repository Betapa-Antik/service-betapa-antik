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

type UbahKataSandiRequest struct {
	KataSandiLama           string `json:"kata_sandi_lama" validate:"required"`
	KataSandiBaru           string `json:"kata_sandi_baru" validate:"required,password"`
	KonfirmasiKataSandiBaru string `json:"konfirmasi_kata_sandi_baru" validate:"required,eqfield=KataSandiBaru"`
}

type LupaKataSandiRequest struct {
	Email       string    `json:"email" validate:"required"`
	PuskesmasID uuid.UUID `json:"puskesmas_id" validate:"required"`
}

type AturUlangKataSandiRequest struct {
	KataSandiBaru           string `json:"kata_sandi_baru" validate:"required,password"`
	KonfirmasiKataSandiBaru string `json:"konfirmasi_kata_sandi_baru" validate:"required,eqfield=KataSandiBaru"`
}

type UpdateStatusLaporan struct {
	Status     string  `json:"status" validate:"required"`
	Keterangan *string `json:"keterangan,omitempty"`
}

type GetAllLaporanRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}
