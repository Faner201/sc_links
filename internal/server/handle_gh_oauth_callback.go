package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Faner201/sc_links/internal/config"
	"github.com/Faner201/sc_links/internal/model"
	"github.com/labstack/echo/v4"
)

type callbackProvider interface {
	GithubAuthCallback(ctx context.Context, sessionCode string) (*model.User, string, error)
}

func HandleGetGithubCallback(cbProvider callbackProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionCode := c.QueryParam("code")
		if sessionCode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "missin code")
		}

		_, jwt, err := cbProvider.GithubAuthCallback(c.Request().Context(), sessionCode)
		if err != nil {
			log.Printf("error handling github auth callback: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		redirectURL := fmt.Sprintf("%s/auth/token.html?token=%s", config.Get().BaseURL, jwt)
		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}
