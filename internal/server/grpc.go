package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/vsespontanno/gochat-grpc/internal/messaging"
	"github.com/vsespontanno/gochat-grpc/internal/models"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	proto.UnimplementedSenderServer
	kp *messaging.KafkaProducer
}

func NewGRPCServer(kp *messaging.KafkaProducer) *GRPCServer {
	return &GRPCServer{kp: kp}
}

func (s *GRPCServer) SendMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	msg := models.Message{
		Sender:    req.GetSender(),
		Recipient: req.GetRecipient(),
		Content:   req.GetContent(),
		Timestamp: time.Now(),
	}
	log.Printf("Received message: %+v", msg)

	err := s.kp.Produce("test", msg)
	if err != nil {
		return nil, status.Error(codes.Aborted, errors.Wrap(err, "failed to produce message").Error())
	}

	response := &proto.MessageResponse{
		Desc: fmt.Sprintf("Message sent successfull. Content: %s", msg.Content),
	}
	return response, nil
}
