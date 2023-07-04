package user

import (
	"context"
	"sync"

	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
)

type inMemory struct {
	m sync.Map
}

func NewInMemory() *inMemory {
	return &inMemory{}
}

func (s *inMemory) CreateOrUpdate(_ context.Context, user model.User) (*model.User, error) {
	s.m.Store(user.GitHubLogin, user)
	return &user, nil
}

func (s *inMemory) Update(_ context.Context, user model.User) error {
	s.m.Store(user.GitHubLogin, user)
	return nil
}

func (s *inMemory) GetByGithubLogin(_ context.Context, login string) (*model.User, error) {
	if user, ok := s.m.Load(login); ok {
		return user.(*model.User), nil
	}
	return nil, errors.ErrNotFound
}

func (s *inMemory) Deactivate(_ context.Context, login string) error {
	if user, ok := s.m.Load(login); ok {
		user.(*model.User).IsActive = false
		s.m.Store(login, user)
		return nil
	}
	return errors.ErrNotFound
}
