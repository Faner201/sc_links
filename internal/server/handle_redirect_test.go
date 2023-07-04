package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Faner201/sc_links/internal/model/dto"
	"github.com/Faner201/sc_links/internal/server"
	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/Faner201/sc_links/internal/storage/shortening"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRedirect(t *testing.T) {
	t.Run("redirect ro origin URL", func(t *testing.T) {
		const (
			url        = "https://www.google.com"
			identifier = "google"
		)
		var (
			redirect = shorten.NewService(shortening.NewInMemory())
			handler  = server.HandleRedirect(redirect)
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodPost, "/"+identifier, nil)
			e        = echo.New()
			c        = e.NewContext(request, recorder)
		)

		c.SetPath("/:identifier")
		c.SetParamNames("identifier")
		c.SetParamValues(identifier)

		_, err := redirect.Shorten(
			context.Background(),
			dto.ShortenInput{
				RawURL:     url,
				Identifier: mo.Some(identifier),
			},
		)
		require.NoError(t, err)

		require.NoError(t, handler(c))
		assert.Equal(t, http.StatusMovedPermanently, recorder.Code)
		assert.Equal(t, url, recorder.Header().Get("Location"))
	})

	t.Run("returns 404 if identifier is not found", func(t *testing.T) {
		const (
			url        = "https://www.google.com"
			identifier = "google"
		)

		var (
			redirect = shorten.NewService(shortening.NewInMemory())
			handler  = server.HandleRedirect(redirect)
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodPost, "/"+identifier, nil)
			e        = echo.New()
			c        = e.NewContext(request, recorder)
		)

		c.SetPath("/:identifier")
		c.SetParamNames("identifier")
		c.SetParamValues(identifier)

		require.Error(t, handler(c))
	})
}
