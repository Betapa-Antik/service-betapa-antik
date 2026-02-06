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
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func ToPuskesmasResponse(puskesmas models.PuskesmasWithTotal) PuskesmasResponse {
	return PuskesmasResponse{
		ID:            puskesmas.ID,
		Foto:          puskesmas.Foto,
		NamaPuskesmas: puskesmas.NamaPuskesmas,
		Kecmatan:      puskesmas.Kecamatan.NamaKecamatan,
		Kelurahan:     puskesmas.Kelurahan.NamaKelurahan,
		Alamat:        puskesmas.Alamat,
		Latitude:      puskesmas.Latitude,
		Longitude:     puskesmas.Longtitude,
		TotalPetugas:  puskesmas.TotalPetugas,
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
