package kelurahancontroller

import (
	kelurahanrequest "betapa-antik-service/internal/dto/request/kelurahan_request"
	kelurahanresponse "betapa-antik-service/internal/dto/response/kelurahan_response"
	kelurahanservice "betapa-antik-service/internal/services/kelurahan_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type KelurahanController struct {
	kelurahanService kelurahanservice.IKelurahanService
}

func NewKelurahanController(kelurahanService kelurahanservice.IKelurahanService) *KelurahanController {
	return &KelurahanController{
		kelurahanService: kelurahanService,
	}
}

func (k *KelurahanController) CreateKelurahan(ctx echo.Context) error {
	var req kelurahanrequest.CreateKelurahanRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err := k.kelurahanService.CreateKelurahan(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Kelurahan berhasil dibuat", nil)
}

func (k *KelurahanController) GetAllKelurahan(ctx echo.Context) error {
	kecamatanId, err := uuid.Parse(ctx.Param("kecamatanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	req := new(kelurahanrequest.GetAllKelurahanRequest)
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	req.KecamatanId = kecamatanId

	data, total, err := k.kelurahanService.GetAllKelurahan(ctx.Request().Context(), *req)
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

	items := make([]kelurahanresponse.KelurahanResponse, len(data))
	for i, v := range data {
		items[i] = kelurahanresponse.ToKelurahanResponse(*v, true)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Kelurahan Berhasil diambil", items, pagination)
}
func (k *KelurahanController) GetKelurahanById(ctx echo.Context) error {
	kelurahanId, err := uuid.Parse(ctx.Param("kelurahanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	kelurahan, err := k.kelurahanService.GetKelurahanById(ctx.Request().Context(), kelurahanId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Kelurahan Berhasil diambil", kelurahanresponse.ToKelurahanResponse(*kelurahan, true))
}

func (k *KelurahanController) UpdateKelurahan(ctx echo.Context) error {
	kelurahanId, err := uuid.Parse(ctx.Param("kelurahanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	var req kelurahanrequest.UpdateKelurahanRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err = k.kelurahanService.UpdateKelurahan(ctx.Request().Context(), kelurahanId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Kelurahan berhasil di update", nil)
}

func (k *KelurahanController) DeleteKelurahan(ctx echo.Context) error {
	kelurahanId, err := uuid.Parse(ctx.Param("kelurahanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	err = k.kelurahanService.DeleteKelurahan(ctx.Request().Context(), kelurahanId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Kelurahan berhasil dihapus", nil)
}
