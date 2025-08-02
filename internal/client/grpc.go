package client

import (
	"context"
	"fmt"

	"github.com/vsespontanno/gochat-grpc/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	pClient  proto.GreeterClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error while making new client: ", err)
		return nil, err
	}
	client := proto.NewGreeterClient(conn)
	return &GRPCClient{Endpoint: endpoint, pClient: client}, nil
}

func (c *GRPCClient) SayHelloAgain(ctx context.Context, in *proto.HelloReq) (*proto.HelloReply, error) {
	resp, err := c.pClient.SayHelloAgain(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
