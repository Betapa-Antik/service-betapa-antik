package surveyresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"
	"fmt"

	"github.com/google/uuid"
)

type SurveyResponse struct {
	ID           uuid.UUID `json:"id"`
	Responden    any       `json:"responden"`
	Petugas      string    `json:"petugas"`
	Tanggal      string    `json:"tanggal"`
	JenisSurvey  string    `json:"jenis_survey"`
	HI           float64   `json:"hi"`
	CI           float64   `json:"ci"`
	BI           float64   `json:"bi"`
	ABJ          float64   `json:"abj"`
	SectionB     any       `json:"section_b"`
	DetailSurvey []any     `json:"detail_survey"`
	FollowUp     any       `json:"follow_up"`
	Gambar       []any     `json:"gambar_urls"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func ToSurveyResponse(survey models.Survey) SurveyResponse {

	// ============================
	// HITUNG INDEX JENTIK PER SURVEY
	// ============================

	var totalWadah int
	var positifWadah int

	if survey.JenisSurvey == models.JenisSurveyJentik {
		for _, item := range survey.Items {
			if item.JumlahTempatAir != nil {
				totalWadah += *item.JumlahTempatAir
			}
			if item.JumlahPositif != nil {
				positifWadah += *item.JumlahPositif
			}
		}
	}

	// CI
	ci := 0.0
	if totalWadah > 0 {
		ci = (float64(positifWadah) / float64(totalWadah)) * 100
	}

	// HI
	hi := 0.0
	if positifWadah > 0 {
		hi = 100
	}

	// BI
	bi := float64(positifWadah) * 100

	// ABJ
	abj := 100 - hi

	// ============================
	// RESPONDEN
	// ============================

	responden := map[string]interface{}{
		"nama_kepala_keluarga": survey.Keluarga.NamaKepalaKeluarga,
		"alamat":               survey.Keluarga.Alamat,
		"kecamatan":            survey.Keluarga.Kecamatan.NamaKecamatan,
		"kelurahan":            survey.Keluarga.Kelurahan.NamaKelurahan,
		"rt":                   survey.Keluarga.RT,
		"rw":                   survey.Keluarga.RW,
	}

	// ============================
	// DETAIL SURVEY
	// ============================

	detailSurvey := []any{}
	for _, item := range survey.Items {

		resp := map[string]interface{}{
			"id":          item.ID,
			"nama_lokasi": item.Lokasi.NamaLokasi,
			"ditemukan":   item.Ditemukan,
			"keterangan":  item.Keterangan,
		}

		if survey.JenisSurvey == models.JenisSurveyNyamuk {
			resp["jumlah_nyamuk"] = item.JumlahNyamuk
			resp["jenis_perkiraan"] = item.JenisPerkiraan
		}

		if survey.JenisSurvey == models.JenisSurveyJentik {
			resp["jumlah_tempat_air"] = item.JumlahTempatAir
			resp["jumlah_positif"] = item.JumlahPositif
		}

		detailSurvey = append(detailSurvey, resp)
	}

	// ============================
	// FOLLOW UP
	// ============================

	var followUp any = nil

	if survey.JenisSurvey == models.JenisSurveyJentik && survey.FollowUpJentik.ID != uuid.Nil {
		followUp = map[string]interface{}{
			"edukasi_psn":   survey.FollowUpJentik.EdukasiPSN,
			"tindak_lanjut": survey.FollowUpJentik.TindakLanjut,
			"catatan":       survey.FollowUpJentik.Catatan,
		}
	}
	if survey.JenisSurvey == models.JenisSurveyNyamuk && survey.FollowUpNyamuk.ID != uuid.Nil {
		followUp = map[string]interface{}{
			"ditemukan_aedes":   survey.FollowUpNyamuk.DitemukanAedes,
			"tingkat_infestasi": survey.FollowUpNyamuk.TingkatInfestasi,
			"fooging_status":    survey.FollowUpNyamuk.FoogingStatus,
			"edukasi_or_abate":  survey.FollowUpNyamuk.EdukasiOrAbate,
			"catatan":           survey.FollowUpNyamuk.Catatan,
		}
	}

	// ============================
	// GAMBAR
	// ============================

	gambars := []any{}
	for _, sg := range survey.SurveyGambar {
		gambars = append(gambars, map[string]interface{}{
			"id":  sg.ID,
			"url": sg.Gambar.Path,
		})
	}

	var sectionB any = nil

	if survey.JenisSurvey == models.JenisSurveyNyamuk && survey.SurveyNyamukInfo.ID != uuid.Nil {
		sectionB = map[string]interface{}{
			"jenis_bangunan":     survey.SurveyNyamukInfo.JenisBangunan,
			"kondisi_lingkungan": survey.SurveyNyamukInfo.KondisiLingkungan,
		}
	}

	if survey.JenisSurvey == models.JenisSurveyJentik && survey.SurveyPSN.ID != uuid.Nil {
		sectionB = map[string]interface{}{
			"menguras_bak_mandi":         survey.SurveyPSN.MengurasBakMandi,
			"menutup_tempat_air":         survey.SurveyPSN.MenutupTempatAir,
			"mendaur_ulang_barang_bekas": survey.SurveyPSN.MendaurUlangBarangBekas,
			"menggunakan_larvasida":      survey.SurveyPSN.MenggunakanLarvasida,
			"menggunakan_kelambu":        survey.SurveyPSN.MenggunakanKelambu,
			"memiliki_tanaman_pengusir":  survey.SurveyPSN.MemilikiTanamanPengusir,
		}
	}

	fmt.Println("Petugas ID:", survey.Petugas.ID)
	fmt.Println("Kecamatan ID:", survey.Keluarga.Kecamatan.ID)

	// ============================
	// FINAL RESPONSE
	// ============================

	return SurveyResponse{
		ID:          survey.ID,
		Responden:   responden,
		Petugas:     survey.Petugas.NamaLengkap,
		Tanggal:     utils.FormatDate(survey.Tanggal),
		JenisSurvey: survey.JenisSurvey,

		HI:  hi,
		CI:  ci,
		BI:  bi,
		ABJ: abj,

		SectionB:     sectionB,
		DetailSurvey: detailSurvey,
		FollowUp:     followUp,
		Gambar:       gambars,
		CreatedAt:    survey.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    survey.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

type LatestSurvey struct {
	ID           uuid.UUID `json:"id"`
	Tanggal      string    `json:"tanggal"`
	JenisSurvey  string    `json:"jenis_survey"`
	NamaKeluarga string    `json:"nama_keluarga"`
	Alamat       string    `json:"alamat"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
}

func ToLatestSurvey(survey models.Survey) LatestSurvey {
	return LatestSurvey{
		ID:           survey.ID,
		Tanggal:      utils.FormatDate(survey.Tanggal),
		JenisSurvey:  survey.JenisSurvey,
		NamaKeluarga: survey.Keluarga.NamaKepalaKeluarga,
		Alamat:       survey.Keluarga.Alamat,
		Status:       survey.Status,
		CreatedAt:    utils.FormatDate(survey.CreatedAt),
	}
}
