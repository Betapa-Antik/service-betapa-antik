package models

import (
	"time"

	"github.com/google/uuid"
)

type Keluarga struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	NamaKepalaKeluarga string    `gorm:"type:varchar(255)" json:"nama_kepala_keluarga"`
	KecamatanID        uuid.UUID `gorm:"type:uuid;not null" json:"kecamatan_id"`
	Kecamatan          Kecamatan `gorm:"foreignKey:KecamatanID;references:ID;constraint:OnDelete:CASCADE;" json:"kecamatan,omitempty"`
	KelurahanID        uuid.UUID `gorm:"type:uuid;not null" json:"kelurahan_id"`
	Kelurahan          Kelurahan `gorm:"foreignKey:KelurahanID;references:ID;constraint:OnDelete:CASCADE;" json:"kelurahan,omitempty"`
	RT                 string    `gorm:"type:varchar(10)" json:"rt"`
	RW                 string    `gorm:"type:varchar(10)" json:"rw"`
	Alamat             string    `gorm:"type:text" json:"alamat"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Keluarga) TableName() string {
	return "keluarga"
}

type SelectKeluarga struct {
	ID                 uuid.UUID `json:"id"`
	NamaKepalaKeluarga string    `json:"nama_kepala_keluarga"`
	Kecamatan          string    `json:"kecamatan"`
	Kelurahan          string    `json:"kelurahan"`
	RT                 string    `json:"rt"`
	RW                 string    `json:"rw"`
	Alamat             string    `json:"alamat"`
}
