package utils

import (
	"betapa-antik-service/configs"
	"betapa-antik-service/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const DefaultTokenTTL = time.Hour * 24

func GenerateToken(u *models.User, ttl time.Duration) (string, int64, error) {
	if ttl == 0 {
		ttl = DefaultTokenTTL
	}
	exp := time.Now().Add(ttl).Unix()
	claims := jwt.MapClaims{
		"sub":       u.ID.String(),
		"email":     u.Email,
		"role_name": u.Role.Nama,
		"exp":       exp,
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(configs.GetJWTSecret()))
	if err != nil {
		return "", 0, err
	}
	return signed, exp, nil
}

// ParseToken validates and parses a JWT token string and returns the subject (user ID)
// and the role name stored in the token claims.
func ParseToken(tokenStr string) (string, string, error) {
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(configs.GetJWTSecret()), nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		sub, _ := claims["sub"].(string)
		roleName, _ := claims["role_name"].(string)
		return sub, roleName, nil
	}
	return "", "", jwt.ErrTokenMalformed
}
