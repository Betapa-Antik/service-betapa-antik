package models

import (
	"time"

	"github.com/google/uuid"
)

type SurveyPSN struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SurveyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"survey_id"`
	Survey   *Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`

	// SECTION B (PRAKTIK PSN 3M PLUS)

	MengurasBakMandi        bool `gorm:"type:boolean;default:false" json:"menguras_bak_mandi"`
	MenutupTempatAir        bool `gorm:"type:boolean;default:false" json:"menutup_tempat_air"`
	MendaurUlangBarangBekas bool `gorm:"type:boolean;default:false" json:"mendaur_ulang_barang_bekas"`
	MenggunakanLarvasida    bool `gorm:"type:boolean;default:false" json:"menggunakan_larvasida"`
	MenggunakanKelambu      bool `gorm:"type:boolean;default:false" json:"menggunakan_kelambu"`
	MemilikiTanamanPengusir bool `gorm:"type:boolean;default:false" json:"memiliki_tanaman_pengusir"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyPSN) TableName() string {
	return "survey_psn"
}
