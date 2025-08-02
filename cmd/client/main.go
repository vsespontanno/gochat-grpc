package main

import (
	"log"

	"github.com/vsespontanno/gochat-grpc/internal/client"
	"github.com/vsespontanno/gochat-grpc/internal/messaging"
)

const endpoint = "localhost:8081"

func main() {
	grpcClient, err := client.NewGRPCClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := messaging.NewKafkaConsumer("localhost:9092", "my-group", *grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Subscribe("test")
}
