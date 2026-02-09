package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ForgotPasswordStatusPending    = "Pending"
	ForgotPasswordStatusPeninjauan = "Peninjauan"
	ForgotPasswordStatusDitolak    = "Ditolak"
	ForgotPasswordStatusDisetujui  = "Disetujui"
)

type LupaKataSandi struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:varchar(255);not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Status    string    `gorm:"varchar(50);not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (LupaKataSandi) TableName() string {
	return "lupa_kata_sandi"
}
