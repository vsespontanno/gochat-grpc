package server

import (
	"context"
	"log"
	"time"

	"github.com/vsespontanno/gochat-grpc/internal/messaging"
	"github.com/vsespontanno/gochat-grpc/internal/models"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
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
	log.Printf("Received message: %+v", msg)

	s.kp.Produce("test", msg)
	return nil, nil
}
