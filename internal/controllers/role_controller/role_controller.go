package rolecontroller

import (
	"net/http"

	rolerequest "betapa-antik-service/internal/dto/request/role_request"
	roleresponse "betapa-antik-service/internal/dto/response/role_response"
	roleservice "betapa-antik-service/internal/services/role_service"
	errormessage "betapa-antik-service/pkg/constant/error_message"
	resp "betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoleController struct {
	roleService roleservice.IRoleService
}

func NewRoleController(roleService roleservice.IRoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (c *RoleController) Create(ctx echo.Context) error {
	var req rolerequest.CreateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err := c.roleService.Create(ctx.Request().Context(), &req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusCreated, "Peran berhasil dibuat", nil)
}

func (c *RoleController) FindAll(ctx echo.Context) error {
	roles, err := c.roleService.FindAll(ctx.Request().Context())
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	var out []roleresponse.RoleResponse
	for _, r := range roles {
		out = append(out, roleresponse.ToRoleResponse(r))
	}
	return resp.Success(ctx, http.StatusOK, "Daftar peran berhasil diambil", out)
}

func (c *RoleController) FindByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("roleId"))
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	role, err := c.roleService.FindByID(ctx.Request().Context(), id)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}

	return resp.Success(ctx, http.StatusOK, "Peran berhasil diambil", roleresponse.ToRoleResponse(*role))
}

func (c *RoleController) Update(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("roleId"))
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	var req rolerequest.UpdateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "Permintaan tidak valid", err.Error())
	}
	if err := ctx.Validate(&req); err != nil {
		validationErrors := utils.ParseValidationError(err)
		return resp.Error(ctx, http.StatusBadRequest, "Validasi gagal", validationErrors)
	}

	err = c.roleService.Update(ctx.Request().Context(), id, &req)
	if err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Peran berhasil diperbarui", nil)
}

func (c *RoleController) Delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("roleId"))
	if err != nil {
		return resp.Error(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}
	if err := c.roleService.Delete(ctx.Request().Context(), id); err != nil {
		if ce, ok := errormessage.AsCustomErr(err); ok {
			return resp.Error(ctx, ce.Status, ce.Msg, ce.Err.Error())
		}
		return resp.Error(ctx, http.StatusInternalServerError, "Terjadi kesalahan", err.Error())
	}
	return resp.Success(ctx, http.StatusOK, "Peran berhasil dihapus", nil)
}
