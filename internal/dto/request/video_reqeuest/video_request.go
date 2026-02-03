package videoreqeuest

type CreateVideoRequest struct {
	Judul     string `form:"judul" validate:"required"`
	Link      string `form:"link" validate:"required"`
	Deskripsi string `form:"deskripsi" validate:"required"`
}

type UpdateVideoRequest struct {
	Judul     string `form:"judul" validate:"required"`
	Link      string `form:"link" validate:"required"`
	Deskripsi string `form:"deskripsi" validate:"required"`
}

type UpdateStatusVideoRequest struct {
	Status string `form:"status" validate:"required"`
}

type GetAllVideoRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}
