package surveyrequest

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateSurveyRequest struct {
	KeluargaID  uuid.UUID `form:"keluarga_id" validate:"required"`
	Tanggal     time.Time `form:"tanggal" validate:"required"`
	JenisSurvey string    `form:"jenis_survey" validate:"required,oneof='Survey Nyamuk' 'Survey Jentik'"`

	// Untuk Survey Nyamuk
	NyamukInfo *CreateSurveyNyamukInfoRequest `form:"nyamuk_info,omitempty"`
	// dikirim sebagai JSON string

	// Untuk Survey Jentik
	PSN *CreateSurveyPSNRequest `form:"survey_psn,omitempty"`
	// dikirim sebagai JSON string
	// ============================
	// Items Lokasi Pengamatan
	// ============================
	Items string `form:"items" validate:"required"`
	// Items dikirim dalam bentuk JSON string (karena multipart)

	// ============================
	// Follow Up Survey (opsional)
	// ============================
	FollowUpNyamuk *CreateSurveyFollowUpNyamukRequest `form:"followup_nyamuk,omitempty"`
	FollowUpJentik *CreateSurveyFollowUpJentikRequest `form:"followup_jentik,omitempty"`

	// ============================
	// Upload Bukti Survey
	// ============================
	Gambar []*multipart.FileHeader `form:"gambar" validate:"omitempty,dive"`
}

type CreateSurveyNyamukInfoRequest struct {
	JenisBangunan     string `json:"jenis_bangunan" validate:"required"`
	KondisiLingkungan string `json:"kondisi_lingkungan" validate:"required"`
}

type CreateSurveyPSNRequest struct {
	MengurasBakMandi        bool `json:"menguras_bak_mandi"`
	MenutupTempatAir        bool `json:"menutup_tempat_air"`
	MendaurUlangBarangBekas bool `json:"mendaur_ulang_barang_bekas"`
	MenggunakanLarvasida    bool `json:"menggunakan_larvasida"`
	MenggunakanKelambu      bool `json:"menggunakan_kelambu"`
	MemilikiTanamanPengusir bool `json:"memiliki_tanaman_pengusir"`
}

type CreateSurveyItemRequest struct {
	NamaLokasi string `json:"nama_lokasi" validate:"required"`
	Ditemukan  bool   `json:"ditemukan"`

	// ============================
	// Untuk Survey Jentik
	// ============================
	JumlahTempatAir *int `json:"jumlah_tempat_air,omitempty"`
	JumlahPositif   *int `json:"jumlah_positif,omitempty"`

	// ============================
	// Untuk Survey Nyamuk Dewasa
	// ============================
	JumlahNyamuk   *int    `json:"jumlah_nyamuk,omitempty"`
	JenisPerkiraan *string `json:"jenis_perkiraan,omitempty"`

	// Catatan umum
	Keterangan *string `json:"keterangan,omitempty"`
}

type CreateSurveyFollowUpNyamukRequest struct {
	DitemukanAedes   bool   `json:"ditemukan_aedes"`
	TingkatInfestasi string `json:"tingkat_infestasi,omitempty"`
	FoogingStatus    string `json:"fooging_status,omitempty"`
	EdukasiOrAbate   bool   `json:"edukasi_or_abate"`
	Catatan          string `json:"catatan,omitempty"`
}

type CreateSurveyFollowUpJentikRequest struct {
	EdukasiPSN   bool   `json:"edukasi_psn"`
	TindakLanjut string `json:"tindak_lanjut,omitempty"`
	Catatan      string `json:"catatan,omitempty"`
}

type GetAllSurveyRequest struct {
	Page        int    `query:"page" validate:"omitempty,min=1"`
	Limit       int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search      string `query:"search" validate:"omitempty"`
	JenisSurvey string `query:"jenis_survey" validate:"omitempty"`
	StartDate   string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate     string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateSurveyRequest struct {
	KeluargaID  uuid.UUID `form:"keluarga_id" validate:"omitempty"`
	Tanggal     time.Time `form:"tanggal" validate:"omitempty"`
	JenisSurvey string    `form:"jenis_survey" validate:"omitempty,oneof='Survey Nyamuk' 'Survey Jentik'"`

	NyamukInfo *CreateSurveyNyamukInfoRequest `form:"nyamuk_info,omitempty"`
	PSN        *CreateSurveyPSNRequest        `form:"survey_psn,omitempty"`
	// ============================
	// Items Lokasi Pengamatan
	// ============================
	Items string `form:"items" validate:"required"`
	// Items dikirim dalam bentuk JSON string (karena multipart)

	// ============================
	// Follow Up Survey (opsional)
	// ============================
	FollowUpNyamuk *CreateSurveyFollowUpNyamukRequest `form:"followup_nyamuk,omitempty"`
	FollowUpJentik *CreateSurveyFollowUpJentikRequest `form:"followup_jentik,omitempty"`

	// ============================
	// Upload Bukti Survey
	// ============================
	GambarBaru     []*multipart.FileHeader `form:"gambar_baru" validate:"omitempty,dive"`
	HapusGambarIDs []uuid.UUID             `form:"hapus_gambar_ids" validate:"omitempty"`
	HapusItemIDs   []uuid.UUID             `form:"hapus_item_ids" validate:"omitempty"`
}

type UpdateSurveyItemRequest struct {
	NamaLokasi string `json:"nama_lokasi,omitempty"`
	Ditemukan  bool   `json:"ditemukan"`

	// ============================
	// Untuk Survey Jentik
	// ============================
	JumlahTempatAir *int `json:"jumlah_tempat_air,omitempty"`
	JumlahPositif   *int `json:"jumlah_positif,omitempty"`

	// ============================
	// Untuk Survey Nyamuk Dewasa
	// ============================
	JumlahNyamuk   *int    `json:"jumlah_nyamuk,omitempty"`
	JenisPerkiraan *string `json:"jenis_perkiraan,omitempty"`

	// Catatan umum
	Keterangan *string `json:"keterangan,omitempty"`
}

type UpdateSurveyFollowUpNyamukRequest struct {
	DitemukanAedes   bool   `json:"ditemukan_aedes,omitempty"`
	TingkatInfestasi string `json:"tingkat_infestasi,omitempty"`
	FoogingStatus    string `json:"fooging_status,omitempty"`
	EdukasiOrAbate   bool   `json:"edukasi_or_abate,omitempty"`
	Catatan          string `json:"catatan,omitempty"`
}

type UpdateSurveyFollowUpJentikRequest struct {
	EdukasiPSN   bool   `json:"edukasi_psn,omitempty"`
	TindakLanjut string `json:"tindak_lanjut,omitempty"`
	Catatan      string `json:"catatan,omitempty"`
}
