package models

import (
	"time"

	"github.com/google/uuid"
)

type Kecamatan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Foto          string    `gorm:"type:varchar(255)" json:"foto"`
	NamaKecamatan string    `gorm:"type:varchar(255);not null" json:"nama_kecamatan"`
	KodeWilayah   string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"kode_wilayah"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Kecamatan) TableName() string {
	return "kecamatan"
}

type KecamatanWithTotal struct {
	Kecamatan
	TotalKelurahan int `json:"total_kelurahan"`
	TotalPuskesmas int `json:"total_puskesmas"`
}
