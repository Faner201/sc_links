package storage

import (
	"context"

	"github.com/Faner201/sc_links/internal/model"
)

type Storage interface {
	Put(ctx context.Context, shortening model.Shortering) (*model.Shortering, error)
	Get(ctx context.Context, identifier string) (*model.Shortering, error)
	IncrementVisits(ctx context.Context, identifier string) error
}
