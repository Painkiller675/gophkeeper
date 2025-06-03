package server

import (
	"context"
	"github.com/Painkiller675/gophkeeper/internal/server/config"
	"github.com/Painkiller675/gophkeeper/internal/server/interceptors"
	"github.com/Painkiller675/gophkeeper/internal/server/services"
)

// Server gRPC server with authentication and user storage services
type Server struct {
	AuthService   *services.AuthService
	SecretService *services.SecretService
	Address       string
}

// New creates a new Server with given settings
func New(cfg config.Config) *Server {
	authService := services.NewAuthService(cfg)
	secretService := services.NewSecretService(cfg)
	return &Server{
		AuthService:   authService,
		SecretService: secretService,
		Address:       cfg.GRPC.Address,
	}
}

// Run launches gRPC server with authentication and user storage services
func (s *Server) Run(ctx context.Context) {
	interceptor := interceptors.NewAuthInterceptor(s.AuthService.TokenManager)

	services.NewServer(
		s.Address,
		services.WithServices(s.AuthService, s.SecretService),
		services.WithUnaryInterceptors(interceptor.Unary()),
		services.WithStreamInterceptors(interceptor.Stream()),
	).Run(ctx)
}
