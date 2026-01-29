package rolerequest

type CreateRoleRequest struct {
	Nama string `json:"nama" validate:"required"`
}

type UpdateRoleRequest struct {
	Nama string `json:"nama" validate:"required"`
}
