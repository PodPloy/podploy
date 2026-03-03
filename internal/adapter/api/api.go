//nolint:revive
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	cmiddleware "github.com/PodPloy/podploy/internal/adapter/api/middlewares"
	"github.com/PodPloy/podploy/internal/domain/ports"
)

// Config holds the parameters needed to create a new HTTP server instance.
type Config struct {
	Host    string
	Port    uint
	Origins []string
	Timeout time.Duration
}

// Server wraps an Echo instance and provides methods for configuring
// middleware, registering routes, and running the HTTP server.
type Server struct {
	server *echo.Echo
	config *Config
	logger ports.ILogger
}

// New creates and returns a Server configured with the given Config and logger.
// It returns an error if cfg or log is nil.
func New(cfg *Config, log ports.ILogger) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	if log == nil {
		return nil, errors.New("log cannot be nil")
	}

	s := &Server{
		server: echo.New(),
		config: cfg,
		logger: log,
	}

	s.setupMiddlewares()
	s.setupBaseRoutes()

	return s, nil
}

func (s *Server) setupBaseRoutes() {
	s.server.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"ok": true})
	})

	s.server.GET("/ready", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"server": true})
	})
}

func (s *Server) setupMiddlewares() {
	s.server.Use(middleware.Recover())
	s.server.Use(middleware.RequestID())
	s.server.Use(cmiddleware.LoggerMiddleware(s.logger))

	s.server.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		HSTSMaxAge:         31536000,
	}))

	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.config.Origins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
}

// Group creates a new route group with the given prefix and optional middleware.
func (s *Server) Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	return s.server.Group(prefix, m...)
}

// RegisterRoute registers a single route with the specified HTTP method, path,
// and handler function. It returns an error if the route cannot be added.
func (s *Server) RegisterRoute(method, path string, handler echo.HandlerFunc) error {
	_, err := s.server.AddRoute(echo.Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
	if err != nil {
		return err
	}

	return nil
}

// Start begins listening for HTTP requests. It blocks until the provided
// context is canceled, at which point it performs a graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	sc := echo.StartConfig{
		Address:         address,
		GracefulTimeout: s.config.Timeout * time.Second,
	}

	s.logger.Info("Start Server HTTP ", ports.String("address", address))

	err := sc.Start(ctx, s.server)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Stop server for error", ports.Error(err))
		return err
	}

	s.logger.Info("Server HTTP stop succefully")

	return nil
}
