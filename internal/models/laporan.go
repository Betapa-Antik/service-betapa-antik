package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	LaporanStatusBaru      = "Baru"
	LaporanStatusPengajuan = "Pengajuan"
	LaporanStatusDisetujui = "Disetujui"
	LaporanStatusDitolak   = "Ditolak"
	LaporanStatusSelesai   = "Selesai"
)

type Laporan struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	NamaPelapor       string          `gorm:"type:varchar(255);not null" json:"nama_pelapor"`
	KontakPelapor     string          `gorm:"type:varchar(255);not null" json:"kontak_pelapor"`
	Alamat            string          `gorm:"type:text;not null" json:"alamat"`
	JudulLaporan      string          `gorm:"type:varchar(255);not null" json:"judul_laporan"`
	DeskripsiLaporan  string          `gorm:"type:text;not null" json:"deskripsi_laporan"`
	PuskesmasID       uuid.UUID       `gorm:"type:uuid;not null" json:"puskesmas_id"`
	Puskesmas         Puskesmas       `gorm:"foreignKey:PuskesmasID;references:ID;constraint:OnDelete:CASCADE;" json:"puskesmas,omitempty"`
	PetugasID         uuid.UUID       `gorm:"type:uuid;default:null" json:"petugas_id,omitempty"`
	Petugas           User            `gorm:"foreignKey:PetugasID;references:ID;constraint:OnDelete:CASCADE;" json:"petugas,omitempty"`
	Status            string          `gorm:"type:varchar(50);not null" json:"status"`
	CatatanAdmin      *string         `gorm:"type:text;default:null" json:"catatan_admin,omitempty"`
	HasilTindakLanjut *string         `gorm:"type:text;default:null" json:"hasil_tindak_lanjut,omitempty"`
	DeliveredAt       *time.Time      `gorm:"default:null" json:"delivered_at,omitempty"`
	DecisionAt        *time.Time      `gorm:"default:null" json:"decision_at,omitempty"`
	FinishedAt        *time.Time      `gorm:"default:null" json:"finished_at,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	LaporanGambar     []LaporanGambar `gorm:"foreignKey:LaporanID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (Laporan) TableName() string {
	return "laporan"
}

type LaporanGambar struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LaporanID uuid.UUID `gorm:"type:uuid;not null" json:"laporan_id"`
	Laporan   Laporan   `gorm:"foreignKey:LaporanID;references:ID;constraint:OnDelete:CASCADE;" json:"laporan,omitempty"`
	GambarID  uuid.UUID `gorm:"type:uuid;not null" json:"gambar_id"`
	Gambar    Gambar    `gorm:"foreignKey:GambarID;references:ID;constraint:OnDelete:CASCADE;" json:"gambar,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (LaporanGambar) TableName() string {
	return "laporan_gambar"
}
