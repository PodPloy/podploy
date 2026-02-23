package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PodPloy/podploy/internal/domain/ports"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	cmiddleware "github.com/PodPloy/podploy/internal/adapter/http/middlewares"
)

type Config struct {
	Host    string
	Port    uint
	Origins []string
}

type Server struct {
	server *echo.Echo
	config *Config
	logger ports.ILogger
}

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

	return s, nil
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

func (s *Server) Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	return s.server.Group(prefix, m...)
}

func (s *Server) RegisterRoute(method, path string, handler echo.HandlerFunc) {
	s.server.AddRoute(echo.Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	sc := echo.StartConfig{
		Address:         address,
		GracefulTimeout: 5 * time.Second,
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
