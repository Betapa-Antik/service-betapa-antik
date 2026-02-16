package models

import "time"

type TotalDataDashboardAdmin struct {
	TotalLaporan         int64   `json:"total_laporan"`
	GrowthPersenLaporan  float64 `json:"growth_persen"`
	TotalKecamatan       int64   `json:"total_kecamatan"`
	TotalPuskesmas       int64   `json:"total_puskesmas"`
	TotalVideo           int64   `json:"total_video"`
	TotalMateri          int64   `json:"total_materi"`
	TotalLaporanPegajuan int64   `json:"total_laporan_pengajuan"`
	TotalLaporanDiterima int64   `json:"total_laporan_diterima"`
	TotalLaporanDitolak  int64   `json:"total_laporan_ditolak"`
	TotalLaporanSelesai  int64   `json:"total_laporan_selesai"`
}

type StatistikDFChart struct {
	Tanggal   time.Time `json:"tanggal"`
	Kecamatan string    `json:"kecamatan"`

	HI  float64 `json:"hi"`
	CI  float64 `json:"ci"`
	BI  float64 `json:"bi"`
	ABJ float64 `json:"abj"`

	DF     float64 `json:"df"`
	Status string  `json:"status"`
}
