package services

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Server a gRPC server
type Server struct {
	Address            string
	Services           []Service
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

// Option defines gRPC service settings
type Option func(*Server)

// WithUnaryInterceptors returns Option, which defines functions-interceptors for some unary RPC requests
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(server *Server) {
		server.UnaryInterceptors = append(server.UnaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors returns Option, which defines function-interceptors for some stream RPC requests
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(server *Server) {
		server.StreamInterceptors = append(server.StreamInterceptors, interceptors...)
	}
}

// WithServices returns Option, which defines services of gRPC server
func WithServices(services ...Service) Option {
	return func(server *Server) {
		server.Services = append(server.Services, services...)
	}
}

// NewServer creates gRPC server with some settings
func NewServer(address string, options ...Option) *Server {
	server := &Server{Address: address}

	for _, option := range options {
		option(server)
	}

	return server
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate & private key
	serverCert, err := tls.LoadX509KeyPair("../../internal/cert/server-cert.pem", "../../internal/cert/server-key.pem")
	if err != nil {
		return nil, err
	}
	// create & return credentials
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

func generateTLSCreds() (credentials.TransportCredentials, error) {
	certFile := "../../cert/server.crt"
	keyFile := "../../cert/server.key"

	return credentials.NewServerTLSFromFile(certFile, keyFile)
}

// Run sets some interceptors, registers services and launches gRPC server
func (s *Server) Run(ctx context.Context) {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start grpc server")
	}

	// get the tls credentials object
	tlsCreds, err := generateTLSCreds()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCreds),
		grpc.ChainUnaryInterceptor(s.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(s.StreamInterceptors...),
	)

	for _, service := range s.Services {
		service.RegisterService(grpcServer)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err = grpcServer.Serve(listen); err != nil {
		log.Error().Err(err).Msg("Error on grpc server Serve")
	}
}
