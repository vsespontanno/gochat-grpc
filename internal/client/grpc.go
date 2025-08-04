package client

import (
	"context"
	"fmt"

	"github.com/vsespontanno/gochat-grpc/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCClient struct {
	Endpoint string
	pClient  proto.SenderClient
	aClient  proto.AuthClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error while making new client: ", err)
		return nil, err
	}
	client := proto.NewSenderClient(conn)
	return &GRPCClient{Endpoint: endpoint, pClient: client}, nil
}

func (c *GRPCClient) SendMessage(ctx context.Context, in *proto.Message) (*proto.None, error) {
	md := metadata.Pairs("authorization", "Bearer your_jwt_token")
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	_, err := c.pClient.SendMessage(ctx, in)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *GRPCClient) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return c.aClient.Register(ctx, in)
}

func (c *GRPCClient) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	return c.aClient.Login(ctx, in)
}
