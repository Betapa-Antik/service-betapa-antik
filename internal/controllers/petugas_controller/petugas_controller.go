package petugascontroller

import (
	petugasrequest "betapa-antik-service/internal/dto/request/petugas_request"
	puskesmasresponse "betapa-antik-service/internal/dto/response/puskesmas_response"
	petugasservice "betapa-antik-service/internal/services/petugas_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"net/http"

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
