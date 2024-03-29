package server

import (
	"context"
	"net/http"

	"github.com/Faner201/sc_links/internal/auth"
	"github.com/Faner201/sc_links/internal/config"
	"github.com/Faner201/sc_links/internal/model/dto"
	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CloserFunc func(context.Context) error

type Server struct {
	e         *echo.Echo
	shortener *shorten.Service
	auth      *auth.Service
	closers   []CloserFunc
}

func New(shortener *shorten.Service, auth *auth.Service) *Server {
	s := &Server{
		shortener: shortener,
		auth:      auth,
	}
	s.setupRouter()

	return s
}

func (s *Server) AddCloser(closer CloserFunc) {
	s.closers = append(s.closers, closer)
}

func (s *Server) setupRouter() {
	s.e = echo.New()
	s.e.HideBanner = true
	s.e.Validator = NewValidator()

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Use(middleware.RequestID())

	s.e.GET("/auth/oauth/github/link", HandleGetGithubAuthLink(s.auth))
	s.e.GET("/auth/oauth/github/callback", HandleGetGithubCallback(s.auth))
	s.e.GET("/auth/token.html", HandleTokenPage())
	s.e.GET("/static/*", HandleStatic())
	restricted := s.e.Group("/api")
	{
		restricted.Use(echojwt.WithConfig(makeJWTConfig(context.Background())))
		restricted.POST("/shorten", HandleShorten(s.shortener))
		restricted.GET("/stats/:identifier", HandleStats(s.shortener))
	}
	s.e.GET("/:identifier", HandleRedirect(s.shortener))

	s.AddCloser(s.e.Shutdown)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.e.ServeHTTP(w, r)
}

func (s *Server) Shutdown(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func makeJWTConfig(ctx context.Context) echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(config.Get().Auth.JWTSecretKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &dto.UserClaims{}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized)
		},
	}
}
