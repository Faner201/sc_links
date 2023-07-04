package shorten

import (
	"context"
	"log"

	"github.com/Faner201/sc_links/internal/model"
	"github.com/Faner201/sc_links/internal/model/dto"
	"github.com/Faner201/sc_links/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(ctx context.Context, input dto.ShortenInput) (*model.Shortering, error) {
	var (
		id         = uuid.New().ID()
		identifier = input.Identifier.OrElse(Shorten(id))
	)

	inputShortening := model.Shortering{
		Identifier:  identifier,
		OriginalURL: input.RawURL,
		CreatedBy:   input.CreatedBy,
	}

	shortening, err := s.storage.Put(ctx, inputShortening)
	if err != nil {
		return nil, err
	}

	return shortening, nil
}

func (s *Service) Get(ctx context.Context, indentifier string) (*model.Shortering, error) {
	shortening, err := s.storage.Get(ctx, indentifier)
	if err != nil {
		return nil, err
	}

	return shortening, nil
}

func (s *Service) Redirect(ctx context.Context, identifier string) (string, error) {
	shortening, err := s.storage.Get(ctx, identifier)
	if err != nil {
		return "", err
	}

	if err := s.storage.IncrementVisits(ctx, identifier); err != nil {
		log.Printf("failed to increment visits for identifier %q: %v", identifier, err)
	}

	return shortening.OriginalURL, nil
}
