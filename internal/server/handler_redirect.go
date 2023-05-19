package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	er "github.com/Faner201/sc_links/internal/error"
	"github.com/labstack/echo/v4"
)

type redirect interface {
	Redirect(ctx context.Context, identifier string) (string, error)
}

func HandleRedirect(redirect redirect) echo.HandlerFunc {
	return func(c echo.Context) error {
		identifier := c.Param("identifier")

		redirectURL, err := redirect.Redirect(c.Request().Context(), identifier)
		if err != nil {
			if errors.Is(err, er.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			log.Printf("error getting redirect url for %q: %v", identifier, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}
