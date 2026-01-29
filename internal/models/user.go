package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Foto        string     `gorm:"type:varchar(255);not null" json:"foto"`
	NamaLengkap string     `gorm:"type:varchar(255);not null" json:"nama_lengkap"`
	NoPegawai   *string    `gorm:"type:varchar(50);nullable;uniqueIndex;not null" json:"no_pegawai"`
	Email       string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PuskesmasID *uuid.UUID `gorm:"type:uuid;nullable;foreignKey:PuskesmasID;references:ID" json:"puskesmas_id"`
	Puskesmas   *Puskesmas `json:"puskesmas,omitempty"`
	RoleID      uuid.UUID  `gorm:"type:uuid;not null;foreignKey:RoleID;references:ID" json:"role_id"`
	Role        Role       `json:"role,omitempty"`
	KataSandi   string     `gorm:"type:varchar(255);not null" json:"kata_sandi"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
