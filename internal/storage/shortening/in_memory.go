package shortening

import (
	"context"
	"sync"
	"time"

	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
)

type inMemory struct {
	m sync.Map
}

func NewInMemory() *inMemory {
	return &inMemory{}
}

func (s *inMemory) Put(_ context.Context, shortening model.Shortering) (*model.Shortering, error) {
	if _, exists := s.m.Load(shortening.Identifier); exists {
		return nil, errors.ErrIdentifiExists
	}

	shortening.CreatedAt = time.Now().UTC()

	s.m.Store(shortening.Identifier, shortening)

	return &shortening, nil
}

func (s *inMemory) Get(_ context.Context, identifier string) (*model.Shortering, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, errors.ErrNotFound
	}

	shortening := v.(model.Shortering)

	return &shortening, nil
}

func (s *inMemory) IncrementVisits(_ context.Context, identifier string) error {
	v, ok := s.m.Load(identifier)
	if !ok {
		return errors.ErrNotFound
	}

	shortening := v.(model.Shortering)

	shortening.Visits++

	s.m.Store(identifier, shortening)

	return nil
}
