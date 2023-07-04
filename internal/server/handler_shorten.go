package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Faner201/sc_links/internal/config"
	er "github.com/Faner201/sc_links/internal/error"
	"github.com/Faner201/sc_links/internal/model"
	"github.com/Faner201/sc_links/internal/model/dto"
	req_dto "github.com/Faner201/sc_links/internal/server/dto"
	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
)

type shortener interface {
	Shorten(ctx context.Context, input dto.ShortenInput) (*model.Shortering, error)
}

func HandleShorten(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req req_dto.ShortenRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(req); err != nil {
			return err
		}

		userToken, ok := c.Get("user").(*jwt.Token)
		if !ok {
			log.Printf("error: user is not presented in context")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		userClaims, ok := userToken.Claims.(*dto.UserClaims)
		if !ok {
			log.Printf("error: failed to get user claims from token")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		identifier := mo.None[string]()
		if strings.TrimSpace(req.Identifier) != "" {
			identifier = mo.Some(req.Identifier)
		}

		input := dto.ShortenInput{
			RawURL:     req.URL,
			Identifier: identifier,
			CreatedBy:  userClaims.User.GitHubLogin,
		}

		shortening, err := shortener.Shorten(c.Request().Context(), input)
		if err != nil {
			var (
				status int
				msg    = err.Error()
			)
			switch {
			case errors.Is(err, er.ErrInvalidURL):
				status = http.StatusBadRequest
			case errors.Is(err, er.ErrIdentifiExists):
				status = http.StatusConflict
			default:
				log.Printf("error shortening url %q: %v", req.URL, err)
				return err
			}
			return c.JSON(status, req_dto.ShortenResponce{Message: msg})
		}

		shortUrl, err := shorten.PrependBaseUrl(config.Get().BaseURL, shortening.Identifier)
		if err != nil {
			log.Printf("error generating full url for %q: %v", shortening.Identifier, err)
			return err
		}

		return c.JSON(
			http.StatusOK,
			req_dto.ShortenResponce{ShortURL: shortUrl},
		)
	}
}
