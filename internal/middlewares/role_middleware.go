package middlewares

import (
	"net/http"
	"strings"

	"betapa-antik-service/configs"
	adminrepo "betapa-antik-service/internal/repositories/admin_repo"
	"betapa-antik-service/pkg/constant/response"
	"betapa-antik-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequireRole returns an Echo middleware that verifies the bearer token, checks it in Redis
// and ensures the user has the required role. It expects an instance of IAdminRepository
// to be passed so it can fetch the user record and attach it to the context.
func RequireRole(role string, repo adminrepo.IAdminRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header")
			}
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized", "Invalid Authorization header")
			}
			token := parts[1]

			// parse token
			userIDStr, roleName, err := utils.ParseToken(token)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized", "Invalid token: "+err.Error())
			}

			// check token in redis
			val, err := configs.GetRedis(c.Request().Context(), "auth:token:"+token)
			if err != nil || val != userIDStr {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized", "Token invalid or expired")
			}

			// parse uuid
			id, err := uuid.Parse(userIDStr)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized", "Invalid token subject")
			}

			// fetch user
			user, err := repo.FindByID(c.Request().Context(), id)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized", "User not found")
			}

			// check role
			if strings.ToLower(roleName) != strings.ToLower(role) && strings.ToLower(user.Role.Nama) != strings.ToLower(role) {
				return response.Error(c, http.StatusForbidden, "Forbidden", "Insufficient role")
			}

			// attach user to context for handlers
			c.Set("user", user)
			return next(c)
		}
	}
}
