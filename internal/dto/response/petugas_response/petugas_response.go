package petugasresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type PetugasPuskesmasResponse struct {
	ID          uuid.UUID `json:"id"`
	Foto        string    `json:"foto"`
	NamaLengkap string    `json:"nama_lengkap"`
	NoPegawai   string    `json:"no_pegawai"`
	Email       string    `json:"email"`
	Puskesmas   string    `json:"puskesmas"`
	Role        string    `json:"jabatan"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToPetugasPuskesmasResponse(petugas models.User) PetugasPuskesmasResponse {
	return PetugasPuskesmasResponse{
		ID:          petugas.ID,
		Foto:        petugas.Foto,
		NamaLengkap: petugas.NamaLengkap,
		NoPegawai:   *petugas.NoPegawai,
		Email:       petugas.Email,
		Puskesmas:   petugas.Puskesmas.NamaPuskesmas,
		Role:        petugas.Role.Nama,
		Status:      petugas.Status,
		CreatedAt:   utils.FormatDate(petugas.CreatedAt),
		UpdatedAt:   utils.FormatDate(petugas.UpdatedAt),
	}
}
