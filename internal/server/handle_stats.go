package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	er "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
	"github.com/labstack/echo/v4"
)

type shorteningProvider interface {
	Get(ctx context.Context, identifier string) (*model.Shortering, error)
}

func HandleStats(provider shorteningProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		identifier := c.Param("identifier")
		shortening, err := provider.Get(c.Request().Context(), identifier)
		if err != nil {
			if errors.Is(err, er.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			log.Printf("failed to get shortening: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get shortening")
		}

		return c.JSON(http.StatusOK, shortening)
	}
}
