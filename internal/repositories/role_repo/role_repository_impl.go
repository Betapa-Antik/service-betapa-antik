package rolerepo

import (
	"betapa-antik-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepositoryImpl(db *gorm.DB) IRoleRepository {
	return &RoleRepositoryImpl{
		db: db,
	}
}

// FindByName implements [IRoleRepository].
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, nama string) (*models.Role, error) {
	var role models.Role

	if err := r.db.WithContext(ctx).First(&role, "nama = ?", nama).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create implements [IRoleRepository].
func (r *RoleRepositoryImpl) Create(ctx context.Context, data *models.Role) error {
	data.ID = uuid.New()
	return r.db.WithContext(ctx).Create(data).Error
}

// FindAll implements [IRoleRepository].
func (r *RoleRepositoryImpl) FindAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// FindByID implements [IRoleRepository].
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update implements [IRoleRepository].
func (r *RoleRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *models.Role) error {
	return r.db.WithContext(ctx).Model(&models.Role{}).Where("id = ?", id).Updates(data).Error
}

// Delete implements [IRoleRepository].
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Role{}).Error
}
