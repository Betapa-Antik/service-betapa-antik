package models

import (
	"time"

	"github.com/google/uuid"
)

type Puskesmas struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Foto          string     `gorm:"type:varchar(255)" json:"foto"`
	NamaPuskesmas string     `gorm:"type:varchar(255);not null" json:"nama_puskesmas"`
	KecamatanID   uuid.UUID  `gorm:"type:uuid;not null;foreignKey:KecamatanID;references:ID" json:"kecamatan_id"`
	Kecamatan     *Kecamatan `json:"kecamatan,omitempty"`
	KelurahanID   *uuid.UUID `gorm:"type:uuid;nullable;foreignKey:KelurahanID;references:ID" json:"kelurahan_id"`
	Kelurahan     *Kelurahan `json:"kelurahan,omitempty"`
	Alamat        string     `gorm:"type:text;not null" json:"alamat"`
	Latitude      float64    `gorm:"type:decimal(10,8)" json:"latitude"`
	Longtitude    float64    `gorm:"type:decimal(11,8)" json:"longtitude"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Puskesmas) TableName() string {
	return "puskesmas"
}
