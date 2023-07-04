package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type githubAuthLinkProvider interface {
	GithubAuthLink() string
}

func HandleGetGithubAuthLink(provider githubAuthLinkProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		link := provider.GithubAuthLink()
		return c.JSON(http.StatusOK, map[string]string{"link": link})
	}
}
