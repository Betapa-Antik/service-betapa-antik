package matericontroller

import (
	materirequest "betapa-antik-service/internal/dto/request/materi_request"
	materiresponse "betapa-antik-service/internal/dto/response/materi_response"
	materiservice "betapa-antik-service/internal/services/materi_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"encoding/json"
	"math"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MateriController struct {
	materiService materiservice.IMateriService
}

func NewMateriController(materiService materiservice.IMateriService) *MateriController {
	return &MateriController{
		materiService: materiService,
	}
}

func (m *MateriController) CreateMateri(ctx echo.Context) error {
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

	var catatan *string
	if v := ctx.FormValue("catatan_tambahan"); v != "" {
		catatan = &v
	}

	req := materirequest.CreateMateriRequest{
		Judul:           ctx.FormValue("judul"),
		Deskripsi:       ctx.FormValue("deskripsi"),
		CatatanTambahan: catatan,
		Gambar:          gambar,
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = m.materiService.CreateMateri(ctx.Request().Context(), &req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Materi berhasil dibuat", nil)
}

func (m *MateriController) GetAllMateri(ctx echo.Context) error {
	req := new(materirequest.GetAllMateriRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := m.materiService.GetAllMateri(ctx.Request().Context(), *req)
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
	items := make([]materiresponse.MateriResponse, len(data))
	for i, value := range data {
		items[i] = materiresponse.ToMateriResponse(*value)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Materi berhasil diambil", items, pagination)
}

func (m *MateriController) GetByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("materiId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	materi, err := m.materiService.GetByID(ctx.Request().Context(), id)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Materi berhasil diambil", materiresponse.ToMateriResponse(*materi))
}

func (m *MateriController) UpdateMateri(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("materiId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	gambar := []*multipart.FileHeader{}
	form, err := ctx.MultipartForm()
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if form != nil {
		if file, ok := form.File["gambar_baru"]; ok {
			gambar = file
		}
	}

	var catatan *string
	if v := ctx.FormValue("catatan_tambahan"); v != "" {
		catatan = &v
	}

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

	req := materirequest.UpdateMateriRequest{
		Judul:           ctx.FormValue("judul"),
		Deskripsi:       ctx.FormValue("deskripsi"),
		CatatanTambahan: catatan,
		GambarBaru:      gambar,
		HapusGambarIDs:  hapusIDs,
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = m.materiService.UpdateMateri(ctx.Request().Context(), id, &req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Materi berhasil diupdate", nil)
}

func (m *MateriController) UpdateStatusMateri(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("materiId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req materirequest.UpdateStatusMateriRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = m.materiService.UpdateStatusMateri(ctx.Request().Context(), id, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Status materi berhasil diupdate", nil)
}

func (m *MateriController) DeleteMateri(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("materiId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	err = m.materiService.DeleteMateri(ctx.Request().Context(), id)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Materi berhasil dihapus", nil)
}
