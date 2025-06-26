package services

import "google.golang.org/grpc"

// Service the interface of the service of the gRPC server
type Service interface {
	RegisterService(grpc.ServiceRegistrar)
}
