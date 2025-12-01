package server

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"

	"github.com/DangeL187/erax"

	"auth/internal/app"
	handler "auth/internal/features/device/handler/grpc"
	pb "auth/internal/infra/grpc/proto/auth"
)

type Server struct {
	grpcServer *grpc.Server
}

func (s *Server) Run(addr string) error {
	zap.L().Info("gRPC server launched on", zap.String("address", addr))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return erax.Wrap(err, "failed to listen on gRPC address")
	}

	if err = s.grpcServer.Serve(lis); err != nil {
		return erax.Wrap(err, "failed to start gRPC server")
	}

	return nil
}

func NewServer(app *app.App) *Server {
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, handler.NewAuthHandler(app))

	return &Server{grpcServer: grpcServer}
}
