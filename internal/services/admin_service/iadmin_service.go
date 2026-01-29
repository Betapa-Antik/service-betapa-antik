package adminservice

import (
	adminrequest "betapa-antik-service/internal/dto/request/admin_request"
	"context"
)

type IAdminService interface {
	Register(ctx context.Context, req *adminrequest.CreateAdminRequest) error
}
