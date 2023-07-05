package dto

import (
	"github.com/Faner201/sc_links/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/mo"
)

type ShortenInput struct {
	RawURL     string
	Identifier mo.Option[string]
	CreatedBy  string
}

type UserClaims struct {
	jwt.RegisteredClaims
	model.User `json:"user_data"`
}
