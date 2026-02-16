package models

type TotalDataDashboardPetugas struct {
	TotalLaporanBaru int64       `json:"laporan_baru"`
	TotalSurvey      int64       `json:"total_survey"`
	TotalLaporan     int64       `json:"total_laporan"`
	TotalKontainer   []Kontainer `json:"total_kontainer"`
}

type Kontainer struct {
	NamaKontainer string `json:"nama_kontainer"`
	Jumlah        int64  `json:"jumlah"`
}
