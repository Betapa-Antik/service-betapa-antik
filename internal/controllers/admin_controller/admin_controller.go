package admincontroller

import (
	"net/http"

	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	adminservice "betapa-antik-service/internal/services/admin_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	resp "betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"

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
