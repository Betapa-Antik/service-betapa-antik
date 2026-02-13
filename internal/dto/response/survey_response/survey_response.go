package surveyresponse

import (
	"betapa-antik-service/internal/models"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
)

type SurveyResponse struct {
	ID           uuid.UUID `json:"id"`
	Responden    any       `json:"responden"`
	Petugas      string    `json:"petugas"`
	Tanggal      string    `json:"tanggal"`
	JenisSurvey  string    `json:"jenis_survey"`
	DetailSurvey []any     `json:"detail_survey"`
	FollowUp     any       `json:"follow_up"`
	Gambar       []string  `json:"gambar_urls"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func ToSurveyResponse(survey models.Survey) SurveyResponse {

	// RESPONDEN
	responden := map[string]interface{}{
		"nama_kepala_keluarga": survey.Keluarga.NamaKepalaKeluarga,
		"alamat":               survey.Keluarga.Alamat,
		"kecamatan":            survey.Keluarga.Kecamatan.NamaKecamatan,
		"kelurahan":            survey.Keluarga.Kelurahan.NamaKelurahan,
		"rt":                   survey.Keluarga.RT,
		"rw":                   survey.Keluarga.RW,
	}

	// PETUGAS
	// petugas := map[string]interface{}{
	// 	"id":   survey.Petugas.ID,
	// 	"nama": survey.Petugas.NamaLengkap,
	// }

	// DETAIL SURVEY
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

	// FOLLOW UP SAFE
	var followUp any = nil

	if survey.JenisSurvey == models.JenisSurveyNyamuk && survey.FollowUpNyamuk.ID != uuid.Nil {
		followUp = map[string]interface{}{
			"ditemukan_aedes":   survey.FollowUpNyamuk.DitemukanAedes,
			"tingkat_infestasi": survey.FollowUpNyamuk.TingkatInfestasi,
			"fogging_status":    survey.FollowUpNyamuk.FoogingStatus,
			"edukasi_or_abate":  survey.FollowUpNyamuk.EdukasiOrAbate,
			"catatan":           survey.FollowUpNyamuk.Catatan,
		}
	}

	if survey.JenisSurvey == models.JenisSurveyJentik && survey.FollowUpJentik.ID != uuid.Nil {
		followUp = map[string]interface{}{
			"edukasi_psn":   survey.FollowUpJentik.EdukasiPSN,
			"tindak_lanjut": survey.FollowUpJentik.TindakLanjut,
			"catatan":       survey.FollowUpJentik.Catatan,
		}
	}

	// GAMBAR URL
	gambars := []string{}
	for _, gambar := range survey.SurveyGambar {
		gambars = append(gambars, gambar.Gambar.Path)
	}

	return SurveyResponse{
		ID:           survey.ID,
		Responden:    responden,
		Petugas:      survey.Petugas.NamaLengkap,
		Tanggal:      utils.FormatDate(survey.Tanggal),
		JenisSurvey:  survey.JenisSurvey,
		DetailSurvey: detailSurvey,
		FollowUp:     followUp,
		Gambar:       gambars,
		CreatedAt:    survey.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    survey.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
