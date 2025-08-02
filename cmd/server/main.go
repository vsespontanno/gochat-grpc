package main

import (
	"fmt"
	"log"
	"net"

	"github.com/vsespontanno/gochat-grpc/internal/messaging"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
	serv "github.com/vsespontanno/gochat-grpc/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcEndpoint := ":8081"
	kp, err := messaging.NewKafkaProducer("localhost:9092")
	if err != nil {
		log.Fatal(err)
	}
	kp.Produce("test", "Hello sakhal doliyl")
	log.Fatal(makeGRPCTransport(grpcEndpoint))
}

func makeGRPCTransport(endpoint string) error {
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	proto.RegisterGreeterServer(server, serv.NewGRPCServer())
	fmt.Println("GRPC transport running on port", endpoint)
	return server.Serve(ln)

}
