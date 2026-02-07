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

// ActiveOrNonActiveAkunPetugas implements [IAdminRepository].
func (a *AdminRepositoryImpl) ActiveOrNonActiveAkunPetugas(ctx context.Context, petugasId uuid.UUID, status string) error {
	return a.db.WithContext(ctx).Model(&models.User{}).Where("id = ? AND status = ?", petugasId, models.UserStatusApproved).Update("status", status).Error
}

// ApproveOrRejectAkunPetugas implements [IAdminRepository].
func (a *AdminRepositoryImpl) ApproveOrRejectAkunPetugas(ctx context.Context, petugasId uuid.UUID, status string) error {
	return a.db.WithContext(ctx).Model(&models.User{}).Where("id = ? AND status = ?", petugasId, models.UserStatusPending).Update("status", status).Error
}

// FindPetugas implements [IAdminRepository].
func (a *AdminRepositoryImpl) FindPetugas(
	ctx context.Context,
	limit int,
	offset int,
	search string,
) ([]*models.User, int, error) {

	var (
		petugas []*models.User
		total   int64
	)

	if limit <= 0 {
		limit = 10
	}

	searchPattern := "%" + search + "%"

	// ============================
	// 1. COUNT QUERY (FAST)
	// ============================
	countQuery := a.db.WithContext(ctx).
		Model(&models.User{}).
		Joins("LEFT JOIN role ON role.id = \"user\".role_id").
		Joins("LEFT JOIN puskesmas ON puskesmas.id = \"user\".puskesmas_id").
		Where("role.nama = ?", "PETUGAS PUSKESMAS")

	if search != "" {
		countQuery = countQuery.Where(`
			"user".nama_lengkap ILIKE ? OR
			puskesmas.nama_puskesmas ILIKE ?
		`, searchPattern, searchPattern)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ============================
	// 2. DATA QUERY
	// ============================
	dataQuery := a.db.WithContext(ctx).
		Model(&models.User{}).
		Preload("Role").
		Preload("Puskesmas").
		Joins("LEFT JOIN role ON role.id = \"user\".role_id").
		Joins("LEFT JOIN puskesmas ON puskesmas.id = \"user\".puskesmas_id").
		Where("role.nama = ?", "PETUGAS PUSKESMAS")

	if search != "" {
		dataQuery = dataQuery.Where(`
			"user".nama_lengkap ILIKE ? OR
			puskesmas.nama_puskesmas ILIKE ?
		`, searchPattern, searchPattern)
	}

	// Pending di atas
	dataQuery = dataQuery.Order(`
		CASE
			WHEN "user".status = 'pending' THEN 1
			WHEN "user".status = 'approved' THEN 3
			WHEN "user".status = 'active' THEN 2
			WHEN "user".status = 'non-active' THEN 4
			WHEN "user".status = 'reject' THEN 5
			ELSE 6
		END
	`)

	// Terbaru
	dataQuery = dataQuery.Order("\"user\".created_at DESC")

	// Pagination
	if err := dataQuery.
		Limit(limit).
		Offset(offset).
		Find(&petugas).Error; err != nil {
		return nil, 0, err
	}

	return petugas, int(total), nil
}
