package gambarrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IGambarRepository interface {
	WithTx(tx *gorm.DB) IGambarRepository
	DeleteByIds(ctx context.Context, ids []uuid.UUID) error
}
