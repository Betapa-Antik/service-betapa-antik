package petugascontroller

import (
	authrequest "betapa-antik-service/internal/dto/request/auth_request"
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	authresponse "betapa-antik-service/internal/dto/response/auth_response"
	petugasresponse "betapa-antik-service/internal/dto/response/petugas_response"
	puskesmasresponse "betapa-antik-service/internal/dto/response/puskesmas_response"
	"betapa-antik-service/internal/models"
	petugasservice "betapa-antik-service/internal/services/petugas_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
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
