package shorten_test

import (
	"context"
	"testing"

	"github.com/Faner201/sc_links/internal/dto"
	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/Faner201/sc_links/internal/storage/shortening"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Shoreten(t *testing.T) {
	t.Run("generates shortening for a given URL", func(t *testing.T) {
		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = dto.ShortenInput{RawURL: "https://www.google.com"}
		)

		shortening, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		assert.NotEmpty(t, shortening.Identifier)
		assert.Equal(t, "https://www.google.com", shortening.OriginalURL)
		assert.NotZero(t, shortening.CreatedAt)
	})

	t.Run("uses custom identifier if provided", func(t *testing.T) {
		const identifier = "google"

		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = dto.ShortenInput{
				RawURL:     "https://www.google.com",
				Identifier: mo.Some(identifier),
			}
		)

		shortening, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		assert.NotEmpty(t, shortening.Identifier)
		assert.Equal(t, "https://www.google.com", shortening.OriginalURL)
		assert.NotZero(t, shortening.CreatedAt)
	})

	t.Run("returns error if identifier is already taken", func(t *testing.T) {
		const identifier = "google"

		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = dto.ShortenInput{
				RawURL:     "https://www.google.com",
				Identifier: mo.Some(identifier),
			}
		)

		_, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		_, err = svc.Shorten(context.Background(), input)
		assert.ErrorIs(t, err, errors.ErrIdentifiExists)
	})
}

func TestService_Redirect(t *testing.T) {
	t.Run("returns redirect URL for a given identifier", func(t *testing.T) {
		const identifier = "google"

		var (
			inMemoryStorage = shortening.NewInMemory()
			svc             = shorten.NewService(inMemoryStorage)
			input           = dto.ShortenInput{
				RawURL:     "https://www.google.com",
				Identifier: mo.Some(identifier),
			}
		)

		shortening, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		redirectURL, err := svc.Redirect(context.Background(), identifier)
		require.NoError(t, err)

		updateShorteting, err := inMemoryStorage.Get(context.Background(), identifier)
		require.NoError(t, err)

		assert.True(t, updateShorteting.Visits-shortening.Visits == 1)
		assert.Equal(t, "https://www.google.com", redirectURL)
	})

	t.Run("return error if identifier if not found", func(t *testing.T) {
		var svc = shorten.NewService(shortening.NewInMemory())

		_, err := svc.Redirect(context.Background(), "google")
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})
}
