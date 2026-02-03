package videocontroller

import (
	videoreqeuest "betapa-antik-service/internal/dto/request/video_reqeuest"
	videoresponse "betapa-antik-service/internal/dto/response/video_response"
	videoservice "betapa-antik-service/internal/services/video_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VideoController struct {
	videoService videoservice.IVideoService
}

func NewVideoController(videoService videoservice.IVideoService) *VideoController {
	return &VideoController{
		videoService: videoService,
	}
}

func (v *VideoController) CreateVideo(ctx echo.Context) error {
	var req videoreqeuest.CreateVideoRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err := v.videoService.CreateVideo(ctx.Request().Context(), req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusCreated, "Video berhasil dibuat", nil)
}

func (v *VideoController) GetAllVideo(ctx echo.Context) error {
	req := new(videoreqeuest.GetAllVideoRequest)
	if err := ctx.Bind(req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}

	data, total, err := v.videoService.GetAllVideo(ctx.Request().Context(), *req)
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
	for i, value := range data {
		items[i] = videoresponse.ToVideoResponse(*value)
	}

	return response.PaginatedSuccess(ctx, http.StatusOK, "Video berhasil diambil", items, pagination)
}

func (v *VideoController) GetVideoById(ctx echo.Context) error {
	videoId, err := uuid.Parse(ctx.Param("videoId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	video, err := v.videoService.GetVideoById(ctx.Request().Context(), videoId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Video berhasil diambil", videoresponse.ToVideoResponse(*video))
}

func (v *VideoController) UpdateVideo(ctx echo.Context) error {
	videoId, err := uuid.Parse(ctx.Param("videoId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req videoreqeuest.UpdateVideoRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err = v.videoService.UpdateVideo(ctx.Request().Context(), videoId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Video berhasil di update", nil)

}

func (v *VideoController) UpdateStatusVideo(ctx echo.Context) error {
	videoId, err := uuid.Parse(ctx.Param("videoId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req videoreqeuest.UpdateStatusVideoRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return response.Error(ctx, http.StatusBadRequest, "validasi gagal", validationErrors)
	}

	err = v.videoService.UpdateStatusVideo(ctx.Request().Context(), videoId, req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Status Video berhasil di update", nil)
}

func (v *VideoController) DeleteVideo(ctx echo.Context) error {
	videoId, err := uuid.Parse(ctx.Param("videoId"))
	if err != nil {
		return response.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	err = v.videoService.DeleteVideo(ctx.Request().Context(), videoId)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return response.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return response.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return response.Success(ctx, http.StatusOK, "Video berhasil di hapus", nil)
}
