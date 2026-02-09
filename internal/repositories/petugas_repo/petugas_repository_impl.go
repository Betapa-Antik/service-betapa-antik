package petugasrepo

import (
	"betapa-antik-service/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PetugasRepositoryImpl struct {
	db *gorm.DB
}

func NewPetugasRepositoryImpl(db *gorm.DB) IPetugasRepository {
	return &PetugasRepositoryImpl{db: db}
}

// DB implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) DB() *gorm.DB {
	return p.db
}

// WithTx implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) WithTx(tx *gorm.DB) IPetugasRepository {
	return NewPetugasRepositoryImpl(tx)
}

// GetSelectPuskesmas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetSelectPuskesmas(ctx context.Context, search string) ([]models.SelectPuskesmas, error) {
	var result []models.SelectPuskesmas

	query := p.db.WithContext(ctx).
		Table("puskesmas").
		Select(`
			puskesmas.id,
			puskesmas.nama_puskesmas,
			kecamatan.nama_kecamatan,
			kelurahan.nama_kelurahan
		`).
		Joins("LEFT JOIN kecamatan ON kecamatan.id = puskesmas.kecamatan_id").
		Joins("LEFT JOIN kelurahan ON kelurahan.id = puskesmas.kelurahan_id").
		Order("puskesmas.nama_puskesmas ASC")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("puskesmas.nama_puskesmas ILIKE ?", searchPattern)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// GetSelectPuskesmasById implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) GetSelectPuskesmasById(ctx context.Context, puskesmasId uuid.UUID) (*models.SelectPuskesmas, error) {
	var puskesmas models.SelectPuskesmas

	query := p.db.WithContext(ctx).
		Table("puskesmas").
		Select(`
			puskesmas.id,
			puskesmas.nama_puskesmas,
			kecamatan.nama_kecamatan,
			kelurahan.nama_kelurahan
		`).
		Joins("LEFT JOIN kecamatan ON kecamatan.id = puskesmas.kecamatan_id").
		Joins("LEFT JOIN kelurahan ON kelurahan.id = puskesmas.kelurahan_id").
		Order("puskesmas.nama_puskesmas ASC")

	if err := query.Where("puskesmas.id = ?", puskesmasId).First(&puskesmas).Error; err != nil {
		return nil, err
	}

	return &puskesmas, nil
}

// RegisterAkunPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) RegisterAkunPetugas(ctx context.Context, data *models.User) error {
	data.ID = uuid.New()
	return p.db.WithContext(ctx).Create(data).Error
}

// UpdateAkunPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) UpdateAkunPetugas(ctx context.Context, petugasId uuid.UUID, updates map[string]interface{}) error {
	return p.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", petugasId).Updates(updates).Error
}

// FindAkunPetugas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindAkunPetugasById(ctx context.Context, petugasId uuid.UUID) (*models.User, error) {
	var petugas models.User
	if err := p.db.WithContext(ctx).Preload("Role").Preload("Puskesmas").Joins("LEFT JOIN role ON role.id = \"user\".role_id").Where("\"user\".id = ? AND role.nama = ?", petugasId, "PETUGAS PUSKESMAS").First(&petugas).Error; err != nil {
		return nil, err
	}

	return &petugas, nil
}

// FindByEmail implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var petugas models.User
	if err := p.db.WithContext(ctx).Preload("Role").Preload("Puskesmas").Joins("LEFT JOIN role ON role.id = \"user\".role_id").Where("\"user\".email = ? AND role.nama = ?", email, "PETUGAS PUSKESMAS").First(&petugas).Error; err != nil {
		return nil, err
	}

	return &petugas, nil
}

// FindLogForgotPasswordByUserID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindLogForgotPasswordByUserID(ctx context.Context, UserId uuid.UUID) (*models.LupaKataSandi, error) {
	var log models.LupaKataSandi

	err := p.db.WithContext(ctx).
		Where("user_id = ? AND status IN ?", UserId,
			[]string{
				string(models.ForgotPasswordStatusPending),
				string(models.ForgotPasswordStatusPeninjauan),
			},
		).
		First(&log).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &log, nil
}

// CreateLogForgotPassword implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) CreateLogForgotPassword(ctx context.Context, data *models.LupaKataSandi) error {
	data.ID = uuid.New()
	return p.db.WithContext(ctx).Create(data).Error
}

// FindPetugasByEmailPuskesmas implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindPetugasByEmailPuskesmas(ctx context.Context, email string, puskesmasId uuid.UUID) (*models.User, error) {
	var user models.User
	if err := p.db.WithContext(ctx).Where("email = ? AND puskesmas_id = ?", email, puskesmasId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindLogForgotPasswordByID implements [IPetugasRepository].
func (p *PetugasRepositoryImpl) FindLogForgotPasswordByID(ctx context.Context, logId uuid.UUID) (*models.LupaKataSandi, error) {
	var log models.LupaKataSandi
	if err := p.db.WithContext(ctx).Where("id = ?", logId).First(&log).Error; err != nil {
		return nil, err
	}

	return &log, nil
}
