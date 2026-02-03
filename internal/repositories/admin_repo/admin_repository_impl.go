package adminrepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepositoryImpl(db *gorm.DB) IAdminRepository {
	return &AdminRepositoryImpl{
		db: db,
	}
}

// Register implements [IAdminRepository].
func (a *AdminRepositoryImpl) Register(ctx context.Context, data *models.User) error {
	data.ID = uuid.New()
	return a.db.WithContext(ctx).Create(data).Error
}

// Update implements [IAdminRepository].
func (a *AdminRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *models.User) error {
	return a.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(data).Error
}

// FindByEmail implements [IAdminRepository].
func (a *AdminRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := a.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID implements [IAdminRepository].
func (a *AdminRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := a.db.WithContext(ctx).Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
