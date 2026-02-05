package puskesmasresponse

import (
	"betapa-antik-service/internal/models"

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
	}
}
