package server

import (
	"context"
	"log"

	"github.com/vsespontanno/gochat-grpc/internal/proto"
)

type GRPCServer struct {
	proto.UnimplementedGreeterServer
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (s *GRPCServer) SayHelloAgain(ctx context.Context, req *proto.HelloReq) (*proto.HelloReply, error) {
	answer := "Hello sakhal doliyl " + req.GetName()
	log.Printf("Received: %v", req.GetName())
	return &proto.HelloReply{Name: answer}, nil
}
