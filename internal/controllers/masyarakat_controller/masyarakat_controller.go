package masyarakatcontroller

import (
	masyarakatrequest "betapa-antik-service/internal/dto/request/masyarakat_request"
	materiresponse "betapa-antik-service/internal/dto/response/materi_response"
	videoresponse "betapa-antik-service/internal/dto/response/video_response"
	masyarakatservice "betapa-antik-service/internal/services/masyarakat_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MasyarakatController struct {
	masyarakatService masyarakatservice.IMasyarakatService
}

func NewMasyarakatController(svc masyarakatservice.IMasyarakatService) *MasyarakatController {
	return &MasyarakatController{
		masyarakatService: svc,
	}
}

func (m *MasyarakatController) GetLatestMaterial(ctx echo.Context) error {
	data, err := m.masyarakatService.GetLatestMaterial(ctx.Request().Context())
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

	return response.Success(ctx, http.StatusOK, "Materi berhasil diambil", items)
}

func (m *MasyarakatController) GetLatestVideo(ctx echo.Context) error {
	data, err := m.masyarakatService.GetLatestVideo(ctx.Request().Context())
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

	return response.Success(ctx, http.StatusOK, "Video berhasil diambil", items)
}

func (m *MasyarakatController) GetAllPublicMateri(ctx echo.Context) error {
	req := new(masyarakatrequest.GetAllPublicMateriRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	data, total, err := m.masyarakatService.GetAllPublicMateri(ctx.Request().Context(), req)
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
	for i, v := range data {
		items[i] = materiresponse.ToMateriResponse(*v)
	}
	return response.PaginatedSuccess(ctx, http.StatusOK, "Materi berhasil diambil", items, pagination)
}

func (m *MasyarakatController) GetAllPublicVideo(ctx echo.Context) error {
	req := new(masyarakatrequest.GetPublicVideoRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	data, total, err := m.masyarakatService.GetAllPublicVideo(ctx.Request().Context(), req)
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
	items := make([]videoresponse.VideoResponse, len(data))
	for i, v := range data {
		items[i] = videoresponse.ToVideoResponse(*v)
	}
	return response.PaginatedSuccess(ctx, http.StatusOK, "Video berhasil diambil", items, pagination)
}

func (m *MasyarakatController) GetPublicMateriByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("materiId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	materi, err := m.masyarakatService.GetPublicMateriByID(ctx.Request().Context(), id)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Materi berhasil diambil", materiresponse.ToMateriResponse(*materi))
}

func (m *MasyarakatController) GetPublicVideoByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("videoId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	video, err := m.masyarakatService.GetPublicVideoByID(ctx.Request().Context(), id)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusOK, "Video berhasil diambil", videoresponse.ToVideoResponse(*video))
}

func (m *MasyarakatController) CreateLaporan(ctx echo.Context) error {
	gambar := []*multipart.FileHeader{}
	form, err := ctx.MultipartForm()
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if form != nil {
		if file, ok := form.File["foto"]; ok {
			gambar = file
		}
	}

	PuskesmasStr := ctx.FormValue("puskesmas_id")
	var PuskesmasID uuid.UUID
	if PuskesmasStr != "" {
		PuskesmasID, err = uuid.Parse(PuskesmasStr)
		if err != nil {
			return response.Error(ctx, 400, "Puskesmas ID tidak valid", err.Error())
		}
	}

	req := masyarakatrequest.CreateLaporanRequest{
		NamaPelapor:      ctx.FormValue("nama_pelapor"),
		KontakPelapor:    ctx.FormValue("kontak_pelapor"),
		Alamat:           ctx.FormValue("alamat"),
		JudulLaporan:     ctx.FormValue("judul_laporan"),
		DeskripsiLaporan: ctx.FormValue("deskripsi_laporan"),
		PuskesmasID:      PuskesmasID,
		Foto:             gambar,
	}

	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = m.masyarakatService.CreateLaporan(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return response.Success(ctx, http.StatusCreated, "Laporan berhasil dibuat", nil)
}

func (m *MasyarakatController) GetLocationResolve(ctx echo.Context) error {
	kelurahan := ctx.QueryParam("kelurahan")
	kecamatan := ctx.QueryParam("kecamatan")

	location, err := m.masyarakatService.GetLocationResolve(ctx.Request().Context(), kelurahan, kecamatan)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Wilayah berhasil diambil", location)
}
