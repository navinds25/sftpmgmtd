package sftpgrpc

import (
	"context"

	pb "github.com/navinds25/sftpmgmt/pkg/sftpevent"
	log "github.com/sirupsen/logrus"
)

// Server is instance of server.
type Server struct{}

// GRPCTest is the function for testing grpc connectivity.
func (s *Server) GRPCTest(ctx context.Context, in *pb.Ack) (*pb.Ack, error) {
	log.Info(in.Message)
	return &pb.Ack{
		Message: "received message.",
	}, nil
}
