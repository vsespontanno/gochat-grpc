package server

import (
	"context"
	"log"
	"time"

	"github.com/vsespontanno/gochat-grpc/internal/messaging"
	"github.com/vsespontanno/gochat-grpc/internal/models"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	proto.UnimplementedSenderServer
	kp *messaging.KafkaProducer
}

func NewGRPCServer(kp *messaging.KafkaProducer) *GRPCServer {
	return &GRPCServer{kp: kp}
}

func (s *GRPCServer) SendMessage(ctx context.Context, req *proto.Message) (*proto.None, error) {
	msg := models.Message{
		Sender:    req.GetSender(),
		Recipient: req.GetRecipient(),
		Content:   req.GetContent(),
		Timestamp: time.Now(),
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	token := md["authorization"]
	log.Printf("Received message: %+v", msg)
	log.Printf("Received message: %s", token)

	s.kp.Produce("test", msg)
	return nil, nil
}
