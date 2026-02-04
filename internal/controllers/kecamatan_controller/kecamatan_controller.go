package kecamatancontroller

import (
	kecamatanrequest "betapa-antik-service/internal/dto/request/kecamatan_request"
	kecamatanresponse "betapa-antik-service/internal/dto/response/kecamatan_response"
	kecamatanservice "betapa-antik-service/internal/services/kecamatan_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type KecamatanController struct {
	kecamatanService kecamatanservice.IKecamatanService
}

func NewKecamatanController(kecamatanService kecamatanservice.IKecamatanService) *KecamatanController {
	return &KecamatanController{kecamatanService: kecamatanService}
}

func (k *KecamatanController) CreateKecamatan(ctx echo.Context) error {
	file, err := ctx.FormFile("foto")
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	req := kecamatanrequest.CreateKecamatanRequest{
		Foto:          file,
		NamaKecamatan: ctx.FormValue("nama_kecamatan"),
		KodeWilayah:   ctx.FormValue("kode_wilayah"),
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err = k.kecamatanService.CreateKecamatan(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Kecamatan Berhasil Dibuat", nil)
}

func (k *KecamatanController) GetAllKecamatan(ctx echo.Context) error {
	req := new(kecamatanrequest.GetAllKecamatanRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := k.kecamatanService.GetAllKecamatan(ctx.Request().Context(), *req)
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

	items := make([]kecamatanresponse.KecamatanResponse, len(data))
	for i, v := range data {
		items[i] = kecamatanresponse.ToKecamatanResponse(*v)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Kecamatan Berhasil diambil", items, pagination)
}

func (k *KecamatanController) GetKecamatanById(ctx echo.Context) error {
	kecmatanId, err := uuid.Parse(ctx.Param("kecamatanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak valid", err.Error())
	}

	kecamatan, err := k.kecamatanService.GetKecamatanById(ctx.Request().Context(), kecmatanId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Kecamatan berhasil diambil", kecamatanresponse.ToKecamatanResponse(*kecamatan))
}

func (k *KecamatanController) UpdateKecamatan(ctx echo.Context) error {
	kecamatanId, err := uuid.Parse(ctx.Param("kecamatanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak valid", err.Error())
	}

	file, err := ctx.FormFile("foto")
	var fotoPtr *multipart.FileHeader

	if err != nil {
		if err == http.ErrMissingFile {
			// Jika file tidak ada, tidak apa-apa (opsional)
			fotoPtr = nil
		} else {
			// Jika error lain (misal: koneksi terputus saat upload)
			return response.Error(ctx, http.StatusBadRequest, "Gagal memproses file", err.Error())
		}
	} else {
		fotoPtr = file
	}

	req := kecamatanrequest.UpdateKecamatanRequest{
		Foto:          fotoPtr, // Field ini harus bertipe *multipart.FileHeader di struct DTO
		NamaKecamatan: ctx.FormValue("nama_kecamatan"),
		KodeWilayah:   ctx.FormValue("kode_wilayah"),
	}

	if err := ctx.Validate(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", utils.ParseValidationError(err))
	}

	err = k.kecamatanService.UpdateKecamatan(ctx.Request().Context(), kecamatanId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Gagal memperbarui kecamatan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Kecamatan Berhasil Diperbarui", nil)
}

func (k *KecamatanController) DeleteKecamatan(ctx echo.Context) error {
	kecamatanId, err := uuid.Parse(ctx.Param("kecamatanId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID Tidak valid", err.Error())
	}

	err = k.kecamatanService.DeleteKecamatan(ctx.Request().Context(), kecamatanId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Kecamatan berhasil dihapus", nil)
}
