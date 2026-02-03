package videorepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRepositoryImpl struct {
	db *gorm.DB
}

func NewVideoRepositoryImpl(db *gorm.DB) IVideoRepository {
	return &VideoRepositoryImpl{
		db: db,
	}
}

// DB implements [IVideoRepository].
func (v *VideoRepositoryImpl) DB() *gorm.DB {
	return v.db
}

// WithTx implements [IVideoRepository].
func (v *VideoRepositoryImpl) WithTx(tx *gorm.DB) IVideoRepository {
	return NewVideoRepositoryImpl(tx)
}

// CreateVideo implements [IVideoRepository].
func (v *VideoRepositoryImpl) CreateVideo(ctx context.Context, data *models.Video) error {
	data.ID = uuid.New()
	return v.db.WithContext(ctx).Create(data).Error
}

// GetAllVideo implements [IVideoRepository].
func (v *VideoRepositoryImpl) GetAllVideo(ctx context.Context, limit int, offset int, search string) ([]*models.Video, int, error) {
	var (
		videoList []*models.Video
		count     int64
	)

	if limit <= 0 {
		limit = 10
	}
	query := v.db.WithContext(ctx).Model(&models.Video{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("judul ILIKE ? OR deskripsi ILIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&videoList).Error
	if err != nil {
		return nil, 0, err
	}
	return videoList, int(count), nil
}

// GetVideoById implements [IVideoRepository].
func (v *VideoRepositoryImpl) GetVideoById(ctx context.Context, videoId uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := v.db.WithContext(ctx).
		Where("id = ?", videoId).
		First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// UpdateVideo implements [IVideoRepository].
func (v *VideoRepositoryImpl) UpdateVideo(ctx context.Context, videoId uuid.UUID, data *models.Video) error {
	return v.db.WithContext(ctx).Model(&models.Video{}).
		Where("id = ?", videoId).
		Updates(data).Error
}

// UpdateStatusVideo implements [IVideoRepository].
func (v *VideoRepositoryImpl) UpdateStatusVideo(ctx context.Context, videoId uuid.UUID, status string) error {
	return v.db.WithContext(ctx).Model(&models.Video{}).Where("id = ?", videoId).Update("status", status).Error
}

// DeleteVideoById implements [IVideoRepository].
func (v *VideoRepositoryImpl) DeleteVideoById(ctx context.Context, videoId uuid.UUID) error {
	return v.db.WithContext(ctx).Where("id = ?", videoId).Delete(&models.Video{}).Error
}
