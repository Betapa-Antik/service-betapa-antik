package keluargacontroller

import (
	keluargarequest "betapa-antik-service/internal/dto/request/keluarga_request"
	kecamatanresponse "betapa-antik-service/internal/dto/response/kecamatan_response"
	keluargaresponse "betapa-antik-service/internal/dto/response/keluarga_response"
	kelurahanresponse "betapa-antik-service/internal/dto/response/kelurahan_response"
	keluargaservice "betapa-antik-service/internal/services/keluarga_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type KeluargaController struct {
	keluargaService keluargaservice.IKeluargaService
}

func NewKeluargaController(keluargaService keluargaservice.IKeluargaService) *KeluargaController {
	return &KeluargaController{keluargaService: keluargaService}
}

func (k *KeluargaController) CreateKeluarga(ctx echo.Context) error {
	var req keluargarequest.CreateKeluargaRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err := k.keluargaService.CreateKeluarga(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Data Keluarga berhasil dibuat", nil)
}

func (k *KeluargaController) GetAllKeluarga(ctx echo.Context) error {
	req := new(keluargarequest.GetAllKeluargaRequest)
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := k.keluargaService.GetAllKeluarga(ctx.Request().Context(), *req)
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

	items := make([]keluargaresponse.KeluargaResponse, len(data))
	for i, v := range data {
		items[i] = keluargaresponse.ToKeluargaResponse(v)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Data keluarga berhasil diambil", items, pagination)
}

func (k *KeluargaController) GetKeluargaById(ctx echo.Context) error {
	keluargaId, err := uuid.Parse(ctx.Param("keluargaId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	keluarga, err := k.keluargaService.GetKeluargaById(ctx.Request().Context(), keluargaId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Data Keluarga berhasil diambil", keluargaresponse.ToKeluargaResponse(*keluarga))
}

func (k *KeluargaController) UpdateKeluarga(ctx echo.Context) error {
	keluargaId, err := uuid.Parse(ctx.Param("keluargaId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req keluargarequest.UpdateKeluargaRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err = k.keluargaService.UpdateKeluarga(ctx.Request().Context(), keluargaId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Data Kelurga berhasil di update", nil)
}

func (k *KeluargaController) DeleteKeluarga(ctx echo.Context) error {
	keluargaId, err := uuid.Parse(ctx.Param("keluargaId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	err = k.keluargaService.DeleteKeluarga(ctx.Request().Context(), keluargaId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Data Keluarga berhasil dihapus", nil)
}

func (p *KeluargaController) GetSelectKecamatan(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	data, err := p.keluargaService.GetSelectKecamatan(ctx.Request().Context(), search)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	items := make([]kecamatanresponse.KecamatanSelectedResponse, len(data))
	for i, v := range data {
		items[i] = kecamatanresponse.ToKecamatanSelectedResponse(v)
	}

	return response.Success(ctx, http.StatusOK, "Kecamatan berhasil diambil", items)
}

func (p *KeluargaController) GetSelectKelurahan(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	kecamatanId, err := uuid.Parse(ctx.Param("kecamatanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak Valid", err.Error())
	}

	data, err := p.keluargaService.GetSelectKelurahan(ctx.Request().Context(), kecamatanId, search)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	items := make([]kelurahanresponse.KelurahanSelectedResponse, len(data))
	for i, v := range data {
		items[i] = kelurahanresponse.ToKelurahanSelectedRespons(v)
	}

	return response.Success(ctx, http.StatusOK, "Kelurahan berhasil diambil", items)
}
