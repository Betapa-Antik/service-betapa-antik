package adminrequest

import "mime/multipart"

type CreateAdminRequest struct {
	Foto               *multipart.FileHeader `form:"foto" validate:"required"`
	NamaLengkap        string                `json:"nama_lengkap" validate:"required"`
	Email              string                `json:"email" validate:"required,email"`
	Password           string                `json:"password" validate:"required"`
	KonfirmasiPassword string                `json:"konfirmasi_password" validate:"required,eqfield=Password"`
}
