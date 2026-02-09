package puskesmascontroller

import (
	puskesmasrequest "betapa-antik-service/internal/dto/request/puskesmas_request"
	kecamatanresponse "betapa-antik-service/internal/dto/response/kecamatan_response"
	kelurahanresponse "betapa-antik-service/internal/dto/response/kelurahan_response"
	puskesmasresponse "betapa-antik-service/internal/dto/response/puskesmas_response"
	puskesmasservice "betapa-antik-service/internal/services/puskesmas_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PuskesmasController struct {
	puskesmasService puskesmasservice.IPuskesmasService
}

func NewPuskesmasController(service puskesmasservice.IPuskesmasService) *PuskesmasController {
	return &PuskesmasController{puskesmasService: service}
}

func (p *PuskesmasController) CreatePuskesmas(ctx echo.Context) error {
	file, err := ctx.FormFile("foto")
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	kecamatanStr := ctx.FormValue("kecamatan_id")
	kelurahanStr := ctx.FormValue("kelurahan_id")

	var kecamatanID uuid.UUID
	var kelurahanID uuid.UUID

	if kecamatanStr != "" {
		kecamatanID, err = uuid.Parse(kecamatanStr)
		if err != nil {
			return response.Error(ctx, 400, "Kecamatan ID tidak valid", err.Error())
		}
	}

	if kelurahanStr != "" {
		kelurahanID, err = uuid.Parse(kelurahanStr)
		if err != nil {
			return response.Error(ctx, 400, "Kelurahan ID tidak valid", err.Error())
		}
	}

	req := puskesmasrequest.CreatePuskesmasRequest{
		Foto:          file,
		NamaPuskesmas: ctx.FormValue("nama_puskesmas"),
		KecamatanID:   kecamatanID,
		KelurahanID:   kelurahanID,
		Alamat:        ctx.FormValue("alamat"),
		Latitude:      ctx.FormValue("latitude"),
		Longtitude:    ctx.FormValue("longitude"),
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err = p.puskesmasService.CreatePuskesmas(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Puskesmas berhasil dibuat", nil)
}

func (p *PuskesmasController) GetAllPuskesmas(ctx echo.Context) error {
	req := new(puskesmasrequest.GetAllKecamatanRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := p.puskesmasService.GetAllPuskesmas(ctx.Request().Context(), *req)
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

	items := make([]puskesmasresponse.PuskesmasResponse, len(data))
	for i, v := range data {
		items[i] = puskesmasresponse.ToPuskesmasResponse(v)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Puskesmas Berhasil Diambil", items, pagination)
}

func (p *PuskesmasController) GetPuskesmasById(ctx echo.Context) error {
	puskesmasId, err := uuid.Parse(ctx.Param("puskesmasId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak Valid", err.Error())
	}

	puskesmas, err := p.puskesmasService.GetPuskesmasById(ctx.Request().Context(), puskesmasId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Puskesmas berhasil di ambil", puskesmasresponse.ToPuskesmasResponse(*puskesmas))
}

func (p *PuskesmasController) UpdatePuskesmas(ctx echo.Context) error {
	puskesmasId, err := uuid.Parse(ctx.Param("puskesmasId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak Valid", err.Error())
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

	kecamatanStr := ctx.FormValue("kecamatan_id")
	kelurahanStr := ctx.FormValue("kelurahan_id")

	var kecamatanID uuid.UUID
	var kelurahanID uuid.UUID

	if kecamatanStr != "" {
		kecamatanID, err = uuid.Parse(kecamatanStr)
		if err != nil {
			return response.Error(ctx, 400, "Kecamatan ID tidak valid", err.Error())
		}
	}

	if kelurahanStr != "" {
		kelurahanID, err = uuid.Parse(kelurahanStr)
		if err != nil {
			return response.Error(ctx, 400, "Kelurahan ID tidak valid", err.Error())
		}
	}

	req := puskesmasrequest.UpdatePuskesmasRequest{
		Foto:          fotoPtr,
		NamaPuskesmas: ctx.FormValue("nama_puskesmas"),
		KecamatanID:   kecamatanID,
		KelurahanID:   kelurahanID,
		Alamat:        ctx.FormValue("alamat"),
		Latitude:      ctx.FormValue("latitude"),
		Longtitude:    ctx.FormValue("longitude"),
	}

	if err := ctx.Validate(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", utils.ParseValidationError(err))
	}

	err = p.puskesmasService.UpdatePuskesmas(ctx.Request().Context(), puskesmasId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Puskesmas berhasil di update", nil)
}

func (p *PuskesmasController) DeletePuskesmas(ctx echo.Context) error {
	puskesmasId, err := uuid.Parse(ctx.Param("puskesmasId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak Valid", err.Error())
	}

	err = p.puskesmasService.DeletePuskesmas(ctx.Request().Context(), puskesmasId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Berhasil menghapus puskesmas", nil)
}

func (p *PuskesmasController) GetSelectKecamatan(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	data, err := p.puskesmasService.GetSelectKecamatan(ctx.Request().Context(), search)
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

func (p *PuskesmasController) GetSelectKelurahan(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	kecamatanId, err := uuid.Parse(ctx.Param("kecamatanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak Valid", err.Error())
	}

	data, err := p.puskesmasService.GetSelectKelurahan(ctx.Request().Context(), kecamatanId, search)
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
