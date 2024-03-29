package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Faner201/sc_links/internal/model"
	"github.com/Faner201/sc_links/internal/model/dto"
	"github.com/Faner201/sc_links/internal/server"
	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/Faner201/sc_links/internal/storage/shortening"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleStats(t *testing.T) {
	t.Run("returns shortening with given identifier", func(t *testing.T) {
		var (
			provider = shorten.NewService(shortening.NewInMemory())
			handler  = server.HandleStats(provider)
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest("GET", "/stats/abc", nil)
			e        = echo.New()
			c        = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		c.SetPath("/stats/:identifier")
		c.SetParamNames("identifier")
		c.SetParamValues("abc")

		_, err := provider.Shorten(
			context.Background(),
			dto.ShortenInput{
				Identifier: mo.Some("abc"),
				RawURL:     "https://google.com",
				CreatedBy:  "user",
			},
		)

		require.NoError(t, err)

		require.NoError(t, handler(c))

		var s model.Shortering
		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&s))

		assert.Equal(t, "abc", s.Identifier)
		assert.Equal(t, "https://google.com", s.OriginalURL)
		assert.Equal(t, "user", s.CreatedBy)
		assert.Equal(t, uint64(0), s.Visits)
	})

	t.Run("return 404 if shortening with givetn identifier dow not exist", func(t *testing.T) {
		var (
			provider = shorten.NewService(shortening.NewInMemory())
			handler  = server.HandleStats(provider)
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest("GET", "/stats/abc", nil)
			e        = echo.New()
			c        = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		c.SetPath("/stats/:identifier")
		c.SetParamNames("identifier")
		c.SetParamValues("abc")

		var httpErr *echo.HTTPError
		require.ErrorAs(t, handler(c), &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})
}
