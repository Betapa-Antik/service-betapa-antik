package adminrequest

import "mime/multipart"

type CreateAdminRequest struct {
	Foto               *multipart.FileHeader `form:"foto" validate:"required"`
	NamaLengkap        string                `json:"nama_lengkap" validate:"required"`
	Email              string                `json:"email" validate:"required,email"`
	Password           string                `json:"password" validate:"required"`
	KonfirmasiPassword string                `json:"konfirmasi_password" validate:"required,eqfield=Password"`
}

// ProfileRequest used for updating profile (email and nama_lengkap)
// kept in admin_request package as requested
type ProfileRequest struct {
	Email       string `json:"email" validate:"required,email"`
	NamaLengkap string `json:"nama_lengkap" validate:"required"`
}
