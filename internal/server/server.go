package server

import (
	"context"
	"net/http"

	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CloserFunc func(context.Context) error

type Server struct {
	e       *echo.Echo
	svc     *shorten.Service
	closers []CloserFunc
}

func New(svc *shorten.Service) *Server {
	s := &Server{svc: svc}
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

	s.e.POST("/shorten", HandleShorten(s.svc))
	s.e.POST("/login", HandleLogin())
	s.e.GET("/:identifier", HandleRedirect(s.svc))

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
