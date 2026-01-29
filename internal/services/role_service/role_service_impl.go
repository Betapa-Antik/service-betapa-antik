package roleservice

import (
	rolerequest "betapa-antik-service/internal/dto/request/role_request"
	"betapa-antik-service/internal/models"
	rolerepo "betapa-antik-service/internal/repositories/role_repo"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	"context"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleServiceImpl struct {
	roleRepo rolerepo.IRoleRepository
}

func NewRoleServiceImpl(roleRepo rolerepo.IRoleRepository) IRoleService {
	return &RoleServiceImpl{
		roleRepo: roleRepo,
	}
}

// Create implements [IRoleService].
func (r *RoleServiceImpl) Create(ctx context.Context, req *rolerequest.CreateRoleRequest) error {
	// check duplicate by name
	if existing, err := r.roleRepo.FindByName(ctx, req.Nama); err == nil && existing != nil {
		return errormessage.NewCustomError(errormessage.ErrExists, "Peran sudah terdaftar", 409)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errormessage.NewCustomError(err, "Gagal memeriksa keberadaan peran", 500)
	}

	role := &models.Role{
		ID:   uuid.New(),
		Nama: req.Nama,
	}

	if err := r.roleRepo.Create(ctx, role); err != nil {
		return errormessage.NewCustomError(err, "Gagal membuat peran", 500)
	}
	return nil
}

// FindAll implements [IRoleService].
func (r *RoleServiceImpl) FindAll(ctx context.Context) ([]models.Role, error) {
	roles, err := r.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, errormessage.NewCustomError(err, "Gagal mengambil daftar peran", 500)
	}
	return roles, nil
}

// FindByID implements [IRoleService].
func (r *RoleServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	role, err := r.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errormessage.NewCustomError(errormessage.ErrNotFound, "Peran tidak ditemukan", 404)
		}
		return nil, errormessage.NewCustomError(err, "Gagal mengambil peran", 500)
	}
	return role, nil
}

// Update implements [IRoleService].
func (r *RoleServiceImpl) Update(ctx context.Context, id uuid.UUID, req *rolerequest.UpdateRoleRequest) error {
	role, err := r.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "Peran tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal mengambil peran", 500)
	}

	if req.Nama != "" {
		role.Nama = req.Nama
	}

	if err := r.roleRepo.Update(ctx, id, role); err != nil {
		return errormessage.NewCustomError(err, "Gagal memperbarui peran", 500)
	}
	return nil
}

// Delete implements [IRoleService].
func (r *RoleServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// ensure exists
	if _, err := r.roleRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errormessage.NewCustomError(errormessage.ErrNotFound, "Peran tidak ditemukan", 404)
		}
		return errormessage.NewCustomError(err, "Gagal mengambil peran", 500)
	}

	if err := r.roleRepo.Delete(ctx, id); err != nil {
		return errormessage.NewCustomError(err, "Gagal menghapus peran", 500)
	}
	return nil
}
