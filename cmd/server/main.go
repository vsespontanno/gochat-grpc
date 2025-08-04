package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/vsespontanno/gochat-grpc/internal/db"
	"github.com/vsespontanno/gochat-grpc/internal/messaging"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
	"github.com/vsespontanno/gochat-grpc/internal/repository/pg"
	serv "github.com/vsespontanno/gochat-grpc/internal/server"
	"github.com/vsespontanno/gochat-grpc/internal/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	grpcEndpoint := ":8081"
	kp, err := messaging.NewKafkaProducer("localhost:9092")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectToPostgres(os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	userStore := pg.NewUserStore(conn)

	authService := auth.NewAuthService(userStore)
	log.Fatal(makeGRPCTransport(grpcEndpoint, kp, authService))

}

func makeGRPCTransport(endpoint string, kp *messaging.KafkaProducer, authService *auth.AuthService) error {
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.UnaryInterceptor(UnaryServerInterceptor))
	proto.RegisterSenderServer(server, serv.NewGRPCServer(kp))
	proto.RegisterAuthServer(server, authService)
	fmt.Println("GRPC transport running on port", endpoint)
	return server.Serve(ln)

}

func UnaryServerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Printf("Received request on method: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Printf("Sending response from method: %s", info.FullMethod)
	return resp, err
}
