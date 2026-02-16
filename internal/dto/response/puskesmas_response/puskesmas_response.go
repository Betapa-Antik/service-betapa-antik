package puskesmasresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type PuskesmasResponse struct {
	ID            uuid.UUID `json:"id"`
	Foto          string    `json:"foto"`
	NamaPuskesmas string    `json:"nama_puskesmas"`
	Kecmatan      string    `json:"kecamatan"`
	Kelurahan     string    `json:"kelurahan"`
	Alamat        string    `json:"alamat"`
	Latitude      string    `json:"latitude"`
	Longitude     string    `json:"longitude"`
	TotalPetugas  int       `json:"total_petugas"`
	DF            float64   `json:"df"`
	Status        string    `json:"status"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func ToPuskesmasResponse(puskesmas models.PuskesmasWithTotal) PuskesmasResponse {
	return PuskesmasResponse{
		ID:            puskesmas.ID,
		Foto:          puskesmas.Foto,
		NamaPuskesmas: puskesmas.NamaPuskesmas,
		Kecmatan:      puskesmas.NamaKecamatan,
		Kelurahan:     puskesmas.NamaKelurahan,
		Alamat:        puskesmas.Alamat,
		Latitude:      puskesmas.Latitude,
		Longitude:     puskesmas.Longtitude,
		TotalPetugas:  puskesmas.TotalPetugas,
		DF:            puskesmas.DF,
		Status:        puskesmas.Status,
		CreatedAt:     utils.FormatDate(puskesmas.CreatedAt),
		UpdatedAt:     utils.FormatDate(puskesmas.UpdatedAt),
	}
}

type PuskesmasSelectedResponse struct {
	ID            uuid.UUID `json:"id"`
	NamaPuskesmas string    `json:"nama_puskesmas"`
	Kecmatan      string    `json:"kecamatan"`
	Kelurahan     string    `json:"kelurahan"`
}

func ToPuskesmasSelectedResponse(puskesmas models.SelectPuskesmas) PuskesmasSelectedResponse {
	return PuskesmasSelectedResponse{
		ID:            puskesmas.ID,
		NamaPuskesmas: puskesmas.NamaPuskesmas,
		Kecmatan:      puskesmas.NamaKecamatan,
		Kelurahan:     puskesmas.NamaKelurahan,
	}
}

type PetugasPuskesmasWithTotalResponse struct {
	ID          uuid.UUID `json:"id"`
	Foto        string    `json:"foto"`
	NamaLengkap string    `json:"nama_lengkap"`
	NoPegawai   string    `json:"no_pegawai"`
	Email       string    `json:"email"`
	Puskesmas   string    `json:"puskesmas"`
	Role        string    `json:"jabatan"`
	Status      string    `json:"status"`
	TotalSurvey int64     `json:"total_survey"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToPetugasPuskesmasWithTotalResponse(petugas models.PetugasWithTotalSurvey) PetugasPuskesmasWithTotalResponse {
	return PetugasPuskesmasWithTotalResponse{
		ID:          petugas.ID,
		Foto:        petugas.Foto,
		NamaLengkap: petugas.NamaLengkap,
		NoPegawai:   petugas.NoPegawai,
		Email:       petugas.Email,
		Puskesmas:   petugas.Puskesmas,
		Role:        petugas.Role,
		Status:      petugas.Status,
		TotalSurvey: petugas.TotalSurvey,
		CreatedAt:   utils.FormatDateString(petugas.CreatedAt),
		UpdatedAt:   utils.FormatDateString(petugas.UpdatedAt),
	}
}
