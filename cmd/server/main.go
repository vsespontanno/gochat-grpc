package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

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

	jwtSecret := os.Getenv("JWT_SECRET")

	jwtService, err := auth.NewJwtService(jwtSecret)
	if err != nil {
		log.Fatal(err)
	}
	userStore := pg.NewUserStore(conn)

	authService := auth.NewAuthService(userStore, jwtService, time.Duration(1)*time.Hour)
	log.Fatal(makeGRPCTransport(grpcEndpoint, kp, authService, jwtService))

}

func makeGRPCTransport(endpoint string, kp *messaging.KafkaProducer, authService *auth.AuthService, jwtService *auth.JwtService) error {
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()) /* , grpc.UnaryInterceptor(NewJWTUnaryInterceptor(jwtService)) */)
	proto.RegisterSenderServer(server, serv.NewGRPCServer(kp))
	proto.RegisterAuthServer(server, authService)
	fmt.Println("GRPC transport running on port", endpoint)
	return server.Serve(ln)

}

// func jwtFromContext(ctx context.Context) (string, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return "", status.Errorf(codes.InvalidArgument, "не удалось получить метаданные")
// 	}
// 	authHeaders, ok := md["authorization"]
// 	if !ok || len(authHeaders) == 0 {
// 		return "", status.Errorf(codes.Unauthenticated, "отсутствует заголовок авторизации")
// 	}
// 	authHeader := authHeaders[0]
// 	return authHeader, nil
// }

// func NewJWTUnaryInterceptor(jwtService *auth.JwtService) grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,

// 	) (interface{}, error) {
// 		method, _ := grpc.Method(ctx)

// 		if method == "/Auth/Login" || method == "/Auth/Register" {
// 			return handler(ctx, req)
// 		}
// 		token, err := jwtFromContext(ctx)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if ok, err := jwtService.ValidateToken(token); !ok || err != nil {
// 			return nil, status.Errorf(codes.Unauthenticated, "недействительный JWT")
// 		}
// 		log.Println("JWT validation success")
// 		return handler(ctx, req)
// 	}
// }
