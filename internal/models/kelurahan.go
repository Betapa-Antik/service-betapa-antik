package models

import (
	"time"

	"github.com/google/uuid"
)

type Kelurahan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	NamaKelurahan string    `gorm:"type:varchar(255);not null" json:"nama_kelurahan"`
	KodeKelurahan string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"kode_kelurahan"`

	KecamatanID uuid.UUID `gorm:"type:uuid;not null" json:"kecamatan_id"`
	Kecamatan   Kecamatan `gorm:"foreignKey:KecamatanID;references:ID;constraint:OnDelete:CASCADE;" json:"kecamatan,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Kelurahan) TableName() string {
	return "kelurahan"
}

type SelectKelurahan struct {
	ID            uuid.UUID `json:"id"`
	NamaKelurahan string    `json:"nama_kelurahan"`
	KodeKelurahan string    `json:"kode_Kelurahan"`
}
