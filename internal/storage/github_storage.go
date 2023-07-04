package storage

import (
	"context"

	"github.com/Faner201/sc_links/internal/model"
)

type StorageGithub interface {
	CreateOrUpdate(context.Context, model.User) (*model.User, error)
	Update(context.Context, model.User) error
	GetByGithubLogin(context.Context, string) (*model.User, error)
	Deactivate(context.Context, string) error
}
