package error

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrIdentifiExists  = errors.New("identifier already exists")
	ErrInvalidURL      = errors.New("invalid url")
	ErrUserIsNotMember = errors.New("user is not member of the organization")
	ErrInvalidToken    = errors.New("invalid token")
	ErrReadConfig      = errors.New("Error loading .env file")
)
