package petugascontroller

import (
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	authresponse "betapa-antik-service/internal/dto/response/auth_response"
	laporanresponse "betapa-antik-service/internal/dto/response/laporan_response"
	logresponse "betapa-antik-service/internal/dto/response/log_response"
	petugasresponse "betapa-antik-service/internal/dto/response/petugas_response"
	puskesmasresponse "betapa-antik-service/internal/dto/response/puskesmas_response"
	surveyresponse "betapa-antik-service/internal/dto/response/survey_response"
	"betapa-antik-service/internal/models"
	petugasservice "betapa-antik-service/internal/services/petugas_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PetugasController struct {
	petugasService petugasservice.IPetugasService
}

func NewPetugasController(petugasService petugasservice.IPetugasService) *PetugasController {
	return &PetugasController{petugasService: petugasService}
}

func (p *PetugasController) GetSelectPuskesmas(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	data, err := p.petugasService.GetSelectPuskesmas(ctx.Request().Context(), search)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	items := make([]puskesmasresponse.PuskesmasSelectedResponse, len(data))
	for i, v := range data {
		items[i] = puskesmasresponse.ToPuskesmasSelectedResponse(v)
	}

	return response.Success(ctx, http.StatusOK, "Puskesmas Berhasil diambil", items)
}

func (p *PetugasController) RegisterAkunPetugas(ctx echo.Context) error {
	file, err := ctx.FormFile("foto")
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	puskemsmasStr := ctx.FormValue("puskesmas_id")
	var puskesmasID uuid.UUID
	if puskemsmasStr != "" {
		puskesmasID, err = uuid.Parse(puskemsmasStr)
		if err != nil {
			return response.Error(ctx, 400, "Puskesmas ID tidak valid", err.Error())
		}
	}

	req := petugasrequest.RegisterPetugasRequest{
		Foto:                file,
		NamaLengkap:         ctx.FormValue("nama_lengkap"),
		NoPegawai:           ctx.FormValue("no_pegawai"),
		Email:               ctx.FormValue("email"),
		PuskesmasId:         puskesmasID,
		KataSandi:           ctx.FormValue("kata_sandi"),
		KonfirmasiKataSandi: ctx.FormValue("konfirmasi_kata_sandi"),
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi Gagal", validationErrors)
	}

	err = p.petugasService.RegisterAkunPetugas(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Akun berhasil dibuat, silahkan tunggu konfirmasi admin", nil)
}

func (p *PetugasController) LoginPetugas(ctx echo.Context) error {
	var req authrequest.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	user, token, err := p.petugasService.LoginPetugas(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	exp := int64(24 * 60 * 60)
	return response.Success(ctx, http.StatusOK, "Login berhasil", authresponse.AuthResponse{Token: token, ExpiresIn: exp, UserID: user.ID})
}

func (p *PetugasController) GetProfilePetugas(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}

	data, err := p.petugasService.GetProfilePetugas(ctx.Request().Context(), user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Berhasil mengambil profil", petugasresponse.ToPetugasPuskesmasResponse(*data))
}

func (p *PetugasController) UpdateProfilePetugas(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}

	file, err := ctx.FormFile("foto")
	var fotoPtr *multipart.FileHeader
	if err != nil {
		if err == http.ErrMissingFile {
			fotoPtr = nil
		} else {
			return response.Error(ctx, http.StatusBadRequest, "Gagal memproses file", err.Error())
		}
	} else {
		fotoPtr = file
	}

	puskemsmasStr := ctx.FormValue("puskesmas_id")
	var puskesmasID uuid.UUID
	if puskemsmasStr != "" {
		puskesmasID, err = uuid.Parse(puskemsmasStr)
		if err != nil {
			return response.Error(ctx, 400, "Puskesmas ID tidak valid", err.Error())
		}
	}

	req := petugasrequest.UpdatePetugasRequest{
		Foto:        fotoPtr,
		NamaLengkap: ctx.FormValue("nama_lengkap"),
		NoPegawai:   ctx.FormValue("no_pegawai"),
		Email:       ctx.FormValue("email"),
		PuskesmasId: puskesmasID,
	}

	if err := ctx.Validate(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", utils.ParseValidationError(err))
	}

	err = p.petugasService.UpdateProfilePetugas(ctx.Request().Context(), user.ID, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Berhasil update profile", nil)
}

func (p *PetugasController) LogoutPetugas(ctx echo.Context) error {
	auth := ctx.Request().Header.Get("Authorization")
	if auth == "" {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header")
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid Authorization header")
	}
	token := parts[1]
	if err := p.petugasService.LogoutPetugas(ctx.Request().Context(), token); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Logout Berhasil", nil)
}

func (p *PetugasController) UbahKataSandiPetugas(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}

	var req petugasrequest.UbahKataSandiRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", 500)
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err := p.petugasService.UbahKataSandi(ctx.Request().Context(), user.ID, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Ubah Kata Sandi berhasil", nil)
}

func (p *PetugasController) LupaKataSandi(ctx echo.Context) error {
	var req petugasrequest.LupaKataSandiRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	LogId, err := p.petugasService.LupaKataSandiRequest(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Permintaan anda sedang diproses oleh admin, mohon ditunggu", LogId)
}

func (p *PetugasController) StatusVerifikasiLupaKataSandi(ctx echo.Context) error {
	logId, err := uuid.Parse(ctx.Param("logId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	log, err := p.petugasService.StatusVerifikasiLupaKataSandi(ctx.Request().Context(), logId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Permintaan lupa kata sandi berhasil diambil", logresponse.ToLogKataSandiResponse(*log))

}

func (p *PetugasController) AturUlangKataSandi(ctx echo.Context) error {
	petugasId, err := uuid.Parse(ctx.Param("petugasId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req petugasrequest.AturUlangKataSandiRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", 500)
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = p.petugasService.AturUlangKataSandi(ctx.Request().Context(), petugasId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Ubah Kata Sandi berhasil", nil)
}

func (p *PetugasController) UpdateStatusLaporan(ctx echo.Context) error {
	laporanId, err := uuid.Parse(ctx.Param("laporanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	var req petugasrequest.UpdateStatusLaporan
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", 500)
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}
	err = p.petugasService.UpdateStatusLaporan(ctx.Request().Context(), laporanId, user.ID, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Status laporan berhasil diperbarui", nil)
}

func (p *PetugasController) GetAllLaporan(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	req := new(petugasrequest.GetAllLaporanRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := p.petugasService.GetAllLaporan(ctx.Request().Context(), *req, user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
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
	items := make([]laporanresponse.LaporanResponse, len(data))
	for i, v := range data {
		items[i] = laporanresponse.ToLaporanResponse(*v)
	}
	return response.PaginatedSuccess(ctx, http.StatusOK, "Laporan berhasil diambil", items, pagination)
}

func (p *PetugasController) GetLaporanByID(ctx echo.Context) error {
	laporanId, err := uuid.Parse(ctx.Param("laporanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	data, err := p.petugasService.GetLaporanByID(ctx.Request().Context(), laporanId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Laporan berhasil diambil", laporanresponse.ToLaporanResponse(*data))
}

func (p *PetugasController) GetDashboard(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	data, err := p.petugasService.GetDashboard(ctx.Request().Context(), user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Dashboard berhasil diambil", data)
}

func (p *PetugasController) GetLatestLaporanByPuskesmasID(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	data, err := p.petugasService.GetLatestLaporanByPuskesmasID(ctx.Request().Context(), user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	items := make([]laporanresponse.LaporanResponse, len(data))
	for i, v := range data {
		items[i] = laporanresponse.ToLaporanResponse(*v)
	}
	return response.Success(ctx, http.StatusOK, "Laporan terbaru berhasil diambil", items)
}

func (p *PetugasController) GetLatestSurveyByPetugasID(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	data, err := p.petugasService.GetLatestSurveyByPetugasID(ctx.Request().Context(), user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	items := make([]surveyresponse.LatestSurvey, len(data))
	for i, v := range data {
		items[i] = surveyresponse.ToLatestSurvey(*v)
	}

	return response.Success(ctx, http.StatusOK, "Survey terbaru berhasil diambil", items)
}
