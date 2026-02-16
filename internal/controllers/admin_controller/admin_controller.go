package admincontroller

import (
	"math"
	"net/http"

	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	adminresponse "betapa-antik-service/internal/dto/response/admin_response"
	authresponse "betapa-antik-service/internal/dto/response/auth_response"
	laporanresponse "betapa-antik-service/internal/dto/response/laporan_response"
	logresponse "betapa-antik-service/internal/dto/response/log_response"
	materiresponse "betapa-antik-service/internal/dto/response/materi_response"
	petugasresponse "betapa-antik-service/internal/dto/response/petugas_response"
	videoresponse "betapa-antik-service/internal/dto/response/video_response"
	"betapa-antik-service/internal/models"
	adminservice "betapa-antik-service/internal/services/admin_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	resp "betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminController struct {
	adminService adminservice.IAdminService
}

func NewAdminController(s adminservice.IAdminService) *AdminController {
	return &AdminController{adminService: s}
}

func (c *AdminController) Register(ctx echo.Context) error {
	// parse multipart form values
	file, err := ctx.FormFile("foto")
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	req := adminrequest.CreateAdminRequest{
		Foto:               file,
		NamaLengkap:        ctx.FormValue("nama_lengkap"),
		Email:              ctx.FormValue("email"),
		Password:           ctx.FormValue("password"),
		KonfirmasiPassword: ctx.FormValue("konfirmasi_password"),
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = c.adminService.Register(ctx.Request().Context(), &req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusCreated, "Admin berhasil dibuat", nil)
}

func (c *AdminController) Login(ctx echo.Context) error {
	var req authrequest.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	user, token, err := c.adminService.Login(ctx.Request().Context(), &req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	// token TTL is default 24h; compute seconds
	exp := int64(24 * 60 * 60)
	return resp.Success(ctx, http.StatusOK, "Login berhasil", authresponse.AuthResponse{Token: token, ExpiresIn: exp, UserID: user.ID})
}

func (c *AdminController) Profile(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return resp.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}
	// get profile from service (cached in redis)
	p, err := c.adminService.GetProfile(ctx.Request().Context(), user.ID)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Berhasil mengambil profil", adminresponse.ToAdminResponse(*p))
}

func (c *AdminController) UpdateProfile(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return resp.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}

	var req adminrequest.ProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	if err := c.adminService.UpdateProfile(ctx.Request().Context(), user.ID, req.Email, req.NamaLengkap); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Profil berhasil diperbarui", nil)
}

func (c *AdminController) UpdateProfilePhoto(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		return resp.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User not found in context")
	}
	user, ok := u.(*models.User)
	if !ok {
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", "Invalid user in context")
	}

	file, err := ctx.FormFile("foto")
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := c.adminService.UpdateProfilePhoto(ctx.Request().Context(), user.ID, file); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Foto profil berhasil diperbarui", nil)
}

func (c *AdminController) Logout(ctx echo.Context) error {
	auth := ctx.Request().Header.Get("Authorization")
	if auth == "" {
		return resp.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header")
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return resp.Error(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid Authorization header")
	}
	token := parts[1]

	if err := c.adminService.Logout(ctx.Request().Context(), token); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Logout berhasil", nil)
}

func (c *AdminController) ApproveOrRejecAkunPetugas(ctx echo.Context) error {
	petugasId, err := uuid.Parse(ctx.Param("petugasId"))
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "ID Tidak valid", err.Error())
	}
	var req adminrequest.UpdateStatusPetugas
	if err := ctx.Bind(&req); err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	if err := c.adminService.ApproveOrRejectAkunPetugas(ctx.Request().Context(), petugasId, req); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusOK, "Berhasil update status petugas", nil)
}

func (c *AdminController) ActiveOrNonActiveAkunPetugas(ctx echo.Context) error {
	petugasId, err := uuid.Parse(ctx.Param("petugasId"))
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "ID Tidak valid", err.Error())
	}
	var req adminrequest.UpdateStatusPetugas
	if err := ctx.Bind(&req); err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	if err := c.adminService.ActiveOrNonActiveAkunPetugas(ctx.Request().Context(), petugasId, req); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusOK, "Berhasil update status petugas", nil)
}

func (c *AdminController) FindPetugas(ctx echo.Context) error {
	req := new(adminrequest.GetAllPetugasRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := c.adminService.FindPetugas(ctx.Request().Context(), *req)
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

	items := make([]petugasresponse.PetugasPuskesmasResponse, len(data))
	for i, v := range data {
		items[i] = petugasresponse.ToPetugasPuskesmasResponse(*v)
	}

	return resp.PaginatedSuccess(ctx, http.StatusOK, "Petugas berhasil diambil", items, pagination)
}

func (c *AdminController) GetActiveLupaKataSandi(ctx echo.Context) error {
	req := new(adminrequest.GetAllLupaKataSandiRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := c.adminService.GetActiveLupaKataSandi(ctx.Request().Context(), *req)
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

	items := make([]logresponse.LupaKataSandiResponse, len(data))
	for i, v := range data {
		items[i] = logresponse.ToLupaKataSandiResponse(*v)
	}

	return resp.PaginatedSuccess(ctx, http.StatusOK, "Permintaan Lupa kata sandi berhasil diambil", items, pagination)

}

func (c *AdminController) UpdateStatusLupaKataSandi(ctx echo.Context) error {
	logId, err := uuid.Parse(ctx.Param("logId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req adminrequest.UpdateStatusPetugas
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = c.adminService.UpdateStatusLupaKataSandi(ctx.Request().Context(), logId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusOK, "Update status berhasil", nil)
}

func (m *AdminController) GetAllLaporan(ctx echo.Context) error {
	req := new(adminrequest.GetAllLaporanRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := m.adminService.GetAllLaporan(ctx.Request().Context(), *req)
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
	items := make([]laporanresponse.LaporanResponse, len(data))
	for i, v := range data {
		items[i] = laporanresponse.ToLaporanResponse(*v)
	}
	return response.PaginatedSuccess(ctx, http.StatusOK, "Laporan berhasil diambil", items, pagination)
}

func (m *AdminController) GetLaporanByID(ctx echo.Context) error {
	laporanId, err := uuid.Parse(ctx.Param("laporanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	data, err := m.adminService.GetLaporanByID(ctx.Request().Context(), laporanId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Laporan berhasil diambil", laporanresponse.ToLaporanResponse(*data))
}

func (m *AdminController) UpdateStatusLaporan(ctx echo.Context) error {
	laporanId, err := uuid.Parse(ctx.Param("laporanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	var req adminrequest.UpdateStatusLaporan
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}
	if err := m.adminService.UpdateStatusLaporan(ctx.Request().Context(), laporanId, req); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Status laporan berhasil diperbarui", nil)
}

func (c *AdminController) GetDashboardAdmin(ctx echo.Context) error {
	data, err := c.adminService.GetDashboardAdmin(ctx.Request().Context())
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusOK, "Dashboard berhasil diambil", data)
}

func (c *AdminController) GetSelectKecamatan(ctx echo.Context) error {
	data, err := c.adminService.GetSelectKecamatan(ctx.Request().Context())
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Kecamatan Berhasil diambil", data)
}

func (c *AdminController) GetStatistikDFChart(ctx echo.Context) error {

	kecamatanIdStr := ctx.QueryParam("kecamatan_id")
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")

	var kecamatanId uuid.UUID
	if kecamatanIdStr != "" {
		kecamatanId, _ = uuid.Parse(kecamatanIdStr)
	}

	data, err := c.adminService.GetStatistikDFChart(
		ctx.Request().Context(),
		kecamatanId,
		startDate,
		endDate,
	)

	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusOK, "Statistik DF Chart berhasil diambil", data)

}

func (c *AdminController) GetLatestMateri(ctx echo.Context) error {
	data, err := c.adminService.GetLatestMateri(ctx.Request().Context())
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	items := make([]materiresponse.MateriResponse, len(data))
	for i, v := range data {
		items[i] = materiresponse.ToMateriResponse(*v)
	}

	return resp.Success(ctx, http.StatusOK, "Materi Berhasil diambil", items)
}

func (c *AdminController) GetLatestVideo(ctx echo.Context) error {
	data, err := c.adminService.GetLatestVideo(ctx.Request().Context())
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	items := make([]videoresponse.VideoResponse, len(data))
	for i, v := range data {
		items[i] = videoresponse.ToVideoResponse(*v)
	}

	return resp.Success(ctx, http.StatusOK, "Video Berhasil diambil", items)
}
