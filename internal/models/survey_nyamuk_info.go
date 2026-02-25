package models

import (
	"time"

	"github.com/google/uuid"
)

type SurveyNyamukInfo struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SurveyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"survey_id"`
	Survey   *Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`

	// SECTION B - INFORMASI UMUM LOKASI

	JenisBangunan     string `gorm:"type:varchar(100);not null" json:"jenis_bangunan"`
	KondisiLingkungan string `gorm:"type:varchar(100);not null" json:"kondisi_lingkungan"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SurveyNyamukInfo) TableName() string {
	return "survey_nyamuk_info"
}
