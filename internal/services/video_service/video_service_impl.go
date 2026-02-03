package videoservice

import (
	"betapa-antik-service/configs"
	videoreqeuest "betapa-antik-service/internal/dto/request/video_reqeuest"
	"betapa-antik-service/internal/models"
	videorepo "betapa-antik-service/internal/repositories/video_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"betapa-antik-service/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VideoServiceImpl struct {
	videoRepo videorepo.IVideoRepository
	rdb       *redis.Client
}

func NewVideoServiceImpl(videoRepo videorepo.IVideoRepository, rdb *redis.Client) IVideoService {
	return &VideoServiceImpl{videoRepo: videoRepo, rdb: rdb}
}

func (v *VideoServiceImpl) InvalidateVideoCache(ctx context.Context, videoId uuid.UUID) {
	_ = configs.DeleteRedis(ctx, "video:"+videoId.String())

	_ = configs.DeleteRedis(ctx, "videos:all:*")
}

// CreateVideo implements [IVideoService].
func (v *VideoServiceImpl) CreateVideo(ctx context.Context, req videoreqeuest.CreateVideoRequest) error {
	video := &models.Video{
		Judul:     req.Judul,
		Link:      req.Link,
		Deskripsi: req.Deskripsi,
		Status:    models.VideoStatusDraft,
	}

	err := utils.RunInTransaction(v.videoRepo.DB(), func(tx *gorm.DB) error {
		repoTx := v.videoRepo.WithTx(tx)
		if err := repoTx.CreateVideo(ctx, video); err != nil {
			return errormessage.NewCustomError(err, "Gagal Membuat Video", 500)
		}
		return nil
	})
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat video", 500)
	}

	v.InvalidateVideoCache(ctx, video.ID)

	return nil
}

// GetAllVideo implements [IVideoService].
func (v *VideoServiceImpl) GetAllVideo(ctx context.Context, req videoreqeuest.GetAllVideoRequest) ([]*models.Video, int, error) {
	key := fmt.Sprintf("videos:all:search:%s:page:%d:limit:%d", req.Search, req.Page, req.Limit)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var (
			videoList []*models.Video
			total     int
		)
		if err := json.Unmarshal([]byte(val), &videoList); err == nil {
			return videoList, total, nil
		}
	}

	page := req.Page
	limit := req.Limit
	search := req.Search
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	data, total, err := v.videoRepo.GetAllVideo(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errormessage.NewCustomError(err, "Gagal mengambil daftar video", 500)
	}
	if len(data) == 0 {
		data = []*models.Video{}
	}

	buf, _ := json.Marshal(map[string]any{
		"videos": data,
		"total":  total,
	})

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)

	return data, total, nil
}

// GetVideoById implements [IVideoService].
func (v *VideoServiceImpl) GetVideoById(ctx context.Context, videoId uuid.UUID) (*models.Video, error) {
	key := fmt.Sprintf("video:%s", videoId)
	val, err := configs.GetRedis(ctx, key)
	if err == nil && val != "" {
		var video models.Video
		if err := json.Unmarshal([]byte(val), &video); err == nil {
			return &video, nil
		}
	}

	video, err := v.videoRepo.GetVideoById(ctx, videoId)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil video", 500)
	}

	buf, _ := json.Marshal(video)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_ = configs.SetRedis(ctx, key, buf, time.Minute*10)
	return video, nil
}

// UpdateVideo implements [IVideoService].
func (v *VideoServiceImpl) UpdateVideo(ctx context.Context, videoId uuid.UUID, req videoreqeuest.UpdateVideoRequest) error {
	return utils.RunInTransaction(v.videoRepo.DB(), func(tx *gorm.DB) error {
		repoTx := v.videoRepo.WithTx(tx)

		video, err := repoTx.GetVideoById(ctx, videoId)
		if err != nil {
			return errormessage.NewCustomError(err, "Gagal mengambil video", 500)
		}

		video.Judul = req.Judul
		video.Link = req.Link
		video.Deskripsi = req.Deskripsi

		if err := repoTx.UpdateVideo(ctx, video.ID, video); err != nil {
			return errormessage.NewCustomError(err, "Gagal mengupdate video", 500)
		}
		v.InvalidateVideoCache(ctx, video.ID)
		return nil
	})
}

// UpdateStatusVideo implements [IVideoService].
func (v *VideoServiceImpl) UpdateStatusVideo(ctx context.Context, videoId uuid.UUID, req videoreqeuest.UpdateStatusVideoRequest) error {
	video, err := v.videoRepo.GetVideoById(ctx, videoId)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mengambil video", 500)
	}
	video.Status = req.Status

	if err := v.videoRepo.UpdateStatusVideo(ctx, video.ID, req.Status); err != nil {
		return errormessage.NewCustomError(err, "Gagal mengupdate status video", 500)
	}

	v.InvalidateVideoCache(ctx, video.ID)
	return nil
}

// DeleteVideo implements [IVideoService].
func (v *VideoServiceImpl) DeleteVideo(ctx context.Context, videoId uuid.UUID) error {
	video, err := v.videoRepo.GetVideoById(ctx, videoId)
	if err != nil {
		return errormessage.NewCustomError(err, "Gagal mengambil video", 500)
	}
	if err := v.videoRepo.DeleteVideoById(ctx, video.ID); err != nil {
		return errormessage.NewCustomError(err, "Gagal menhapus video", 500)
	}
	v.InvalidateVideoCache(ctx, video.ID)
	return nil
}
