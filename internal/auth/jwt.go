package auth

import (
	"time"

	"github.com/Faner201/sc_links/internal/config"
	"github.com/Faner201/sc_links/internal/model"
	"github.com/Faner201/sc_links/internal/model/dto"
	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(user model.User) (string, error) {
	claims := dto.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "Faner201",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get().Auth.JWTSecretKey))
}
