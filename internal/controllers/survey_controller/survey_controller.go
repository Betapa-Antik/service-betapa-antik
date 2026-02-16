package surveycontroller

import (
	surveyrequest "betapa-antik-service/internal/dto/request/survey_request"
	keluargaresponse "betapa-antik-service/internal/dto/response/keluarga_response"
	surveyresponse "betapa-antik-service/internal/dto/response/survey_response"
	"betapa-antik-service/internal/models"
	surveyservice "betapa-antik-service/internal/services/survey_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"encoding/json"
	"math"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SurveyController struct {
	surveyService surveyservice.ISurveyService
}

func NewSurveyController(surveyService surveyservice.ISurveyService) *SurveyController {
	return &SurveyController{surveyService: surveyService}
}

func (s *SurveyController) GetSelectKeluarga(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	data, err := s.surveyService.GetSelectKeluarga(ctx.Request().Context(), search)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	items := make([]keluargaresponse.KeluargaSelectedResponse, len(data))
	for i, v := range data {
		items[i] = keluargaresponse.ToKeluargaSelectedResponse(v)
	}

	return response.Success(ctx, http.StatusOK, "Data keluarga berhasil di ambil", items)
}

func (s *SurveyController) CreateSurvey(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}

	gambar := []*multipart.FileHeader{}
	form, err := ctx.MultipartForm()
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if form != nil {
		if file, ok := form.File["gambar"]; ok {
			gambar = file
		}
	}

	keluargaStr := ctx.FormValue("keluarga_id")
	var keluargaID uuid.UUID
	if keluargaStr != "" {
		keluargaID, err = uuid.Parse(keluargaStr)
		if err != nil {
			return response.Error(ctx, 400, "Keluarga ID tidak valid", err.Error())
		}
	}

	tanggalStr := ctx.FormValue("tanggal")
	tanggal, err := utils.ParseTanggal(tanggalStr)
	if err != nil {
		return response.Error(ctx, 400,
			"Tanggal tidak valid", err.Error())
	}

	jenisSurvey := ctx.FormValue("jenis_survey")
	if jenisSurvey != models.JenisSurveyNyamuk &&
		jenisSurvey != models.JenisSurveyJentik {
		return response.Error(ctx, 400,
			"Jenis survey tidak valid",
			"Jenis survey harus Survey Nyamuk atau Survey Jentik")
	}

	itemsStr := ctx.FormValue("items")
	if itemsStr == "" {
		return response.Error(ctx, 400,
			"Items wajib diisi", "items kosong")
	}

	var followNyamuk *surveyrequest.CreateSurveyFollowUpNyamukRequest
	var followJentik *surveyrequest.CreateSurveyFollowUpJentikRequest

	if jenisSurvey == models.JenisSurveyNyamuk {

		followStr := ctx.FormValue("followup_nyamuk")
		if followStr == "" {
			return response.Error(ctx, 400,
				"Follow up nyamuk wajib diisi", "")
		}

		var temp surveyrequest.CreateSurveyFollowUpNyamukRequest
		if err := utils.DecodeJSON(followStr, &temp); err != nil {
			return response.Error(ctx, 400,
				"Format followup_nyamuk tidak valid", err.Error())
		}

		followNyamuk = &temp
	}

	if jenisSurvey == models.JenisSurveyJentik {

		followStr := ctx.FormValue("followup_jentik")
		if followStr == "" {
			return response.Error(ctx, 400,
				"Follow up jentik wajib diisi", "")
		}

		var temp surveyrequest.CreateSurveyFollowUpJentikRequest
		if err := utils.DecodeJSON(followStr, &temp); err != nil {
			return response.Error(ctx, 400,
				"Format followup_jentik tidak valid", err.Error())
		}

		followJentik = &temp
	}

	req := surveyrequest.CreateSurveyRequest{
		KeluargaID:     keluargaID,
		Tanggal:        tanggal,
		JenisSurvey:    jenisSurvey,
		Items:          itemsStr,
		FollowUpNyamuk: followNyamuk,
		FollowUpJentik: followJentik,
		Gambar:         gambar,
	}

	err = s.surveyService.CreateSurvey(
		ctx.Request().Context(),
		user.ID,
		req,
	)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, 500,
			"Gagal membuat survey", err.Error())
	}

	return response.Success(ctx, http.StatusCreated,
		"Survey berhasil dibuat", nil)
}

func (s *SurveyController) GetAllSurvey(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	req := new(surveyrequest.GetAllSurveyRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := s.surveyService.GetAllSurvey(ctx.Request().Context(), *req, user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	pagination := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalData:   int(total),
		TotalPages:  totalPages,
	}
	items := make([]surveyresponse.SurveyResponse, len(data))
	for i, v := range data {
		items[i] = surveyresponse.ToSurveyResponse(v)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Survey berhasil diambil", items, pagination)
}

func (s *SurveyController) GetSurveyByID(ctx echo.Context) error {
	surveyID, err := uuid.Parse(ctx.Param("surveyId"))
	if err != nil {
		return response.Error(ctx, 400, "ID survey tidak valid", err.Error())
	}

	data, err := s.surveyService.GetSurveyByID(ctx.Request().Context(), surveyID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, 500,
			"Gagal mengambil data survey", err.Error())
	}

	return response.Success(ctx, http.StatusOK,
		"Survey berhasil diambil", surveyresponse.ToSurveyResponse(*data))
}

func (s *SurveyController) UpdateSurvey(ctx echo.Context) error {
	// 1. Ambil ID dari Parameter
	surveyID, err := uuid.Parse(ctx.Param("surveyId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID survey tidak valid", err.Error())
	}

	// 2. Handling Multipart Form (Gambar Baru & Hapus Gambar)
	var gambarBaru []*multipart.FileHeader
	form, err := ctx.MultipartForm()
	if err == nil && form != nil {
		if files, ok := form.File["gambar_baru"]; ok {
			gambarBaru = files
		}
	}

	// Parsing ID gambar yang ingin dihapus (dikirim sebagai slice of strings)
	var hapusIDs []uuid.UUID
	hapusStr := ctx.FormValue("hapus_gambar_ids")

	if hapusStr != "" {
		if err := json.Unmarshal([]byte(hapusStr), &hapusIDs); err != nil {
			return response.Error(
				ctx,
				http.StatusBadRequest,
				"Format hapus_gambar_ids tidak valid",
				err.Error(),
			)
		}
	}

	// 3. Parsing Data Utama Survey
	keluargaStr := ctx.FormValue("keluarga_id")
	var keluargaID uuid.UUID
	if keluargaStr != "" {
		keluargaID, _ = uuid.Parse(keluargaStr)
	}

	tanggalStr := ctx.FormValue("tanggal")
	var tanggal time.Time
	if tanggalStr != "" {
		tanggal, _ = utils.ParseTanggal(tanggalStr)
	}

	jenisSurvey := ctx.FormValue("jenis_survey")
	itemsStr := ctx.FormValue("items")

	// 4. Parsing Follow Up (Selective)
	var followNyamuk *surveyrequest.CreateSurveyFollowUpNyamukRequest
	var followJentik *surveyrequest.CreateSurveyFollowUpJentikRequest

	if jenisSurvey == models.JenisSurveyNyamuk {
		followStr := ctx.FormValue("followup_nyamuk")
		if followStr != "" {
			var temp surveyrequest.CreateSurveyFollowUpNyamukRequest
			if err := utils.DecodeJSON(followStr, &temp); err == nil {
				followNyamuk = &temp
			}
		}
	}

	if jenisSurvey == models.JenisSurveyJentik {
		followStr := ctx.FormValue("followup_jentik")
		if followStr != "" {
			var temp surveyrequest.CreateSurveyFollowUpJentikRequest
			if err := utils.DecodeJSON(followStr, &temp); err == nil {
				followJentik = &temp
			}
		}
	}
	// Di dalam UpdateSurvey controller
	var hapusItemIDs []uuid.UUID
	hapusItemStr := ctx.FormValue("hapus_item_ids")

	if hapusItemStr != "" {
		if err := json.Unmarshal([]byte(hapusItemStr), &hapusItemIDs); err != nil {
			return response.Error(
				ctx,
				http.StatusBadRequest,
				"Format hapus_item_ids tidak valid",
				err.Error(),
			)
		}
	}

	// 5. Susun Request DTO
	req := surveyrequest.UpdateSurveyRequest{
		KeluargaID:     keluargaID,
		Tanggal:        tanggal,
		JenisSurvey:    jenisSurvey,
		Items:          itemsStr,
		FollowUpNyamuk: followNyamuk,
		FollowUpJentik: followJentik,
		GambarBaru:     gambarBaru,
		HapusGambarIDs: hapusIDs,
		HapusItemIDs:   hapusItemIDs,
	}

	// 6. Panggil Service
	err = s.surveyService.UpdateSurvey(ctx.Request().Context(), surveyID, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Gagal mengupdate survey", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Survey berhasil diperbarui", nil)
}

func (s *SurveyController) DeleteSurvey(ctx echo.Context) error {
	surveyID, err := uuid.Parse(ctx.Param("surveyId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID survey tidak valid", err.Error())
	}
	err = s.surveyService.DeleteSurvey(ctx.Request().Context(), surveyID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Gagal menghapus survey", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Survey berhasil dihapus", nil)
}
