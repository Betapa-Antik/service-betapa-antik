package models

import (
	"time"

	"github.com/google/uuid"
)

type SurveyItem struct {
	ID              uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	SurveyID        uuid.UUID    `gorm:"type:uuid;not null" json:"survey_id"`
	Survey          *Survey      `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`
	LokasiID        uuid.UUID    `gorm:"not null" json:"lokasi_id"`
	Lokasi          SurveyLokasi `gorm:"foreignKey:LokasiID;references:ID;constraint:OnDelete:CASCADE;" json:"lokasi,omitempty"`
	Ditemukan       bool         `gorm:"type:bool;default:false" json:"ditemukan"`
	JumlahTempatAir *int         `gorm:"type:int" json:"jumlah_tempat_air,omitempty"`
	JumlahNyamuk    *int         `gorm:"type:int" json:"jumlah_nyamuk,omitempty"`
	JumlahPositif   *int         `gorm:"type:int" json:"jumlah_positif,omitempty"`
	JenisPerkiraan  *string      `gorm:"type:varchar(50)" json:"jenis_perkiraan,omitempty"`
	Keterangan      *string      `gorm:"type:text" json:"keterangan,omitempty"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyItem) TableName() string {
	return "survey_item"
}
