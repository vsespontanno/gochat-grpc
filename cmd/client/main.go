package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vsespontanno/gochat-grpc/internal/client"
	"github.com/vsespontanno/gochat-grpc/internal/messaging"
)

const endpoint = "localhost:8081"

type Command string

const (
	HELP     Command = "h"
	WRITE    Command = "w"
	REGISTER Command = "reg"
	LOGIN    Command = "log"
)

func main() {

	ctx := context.Background()
	grpcClient, err := client.NewGRPCClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Connected to server, choose command:\nreg - register\nlog - login\n")
	var c Command
	_, err = fmt.Scan(&c)
	if err != nil {
		log.Fatal("ой")
	}
	if c == REGISTER {
		for {
			err = client.Register(ctx, grpcClient)
			if err == nil {
				break
			}
		}

	}
	userID := client.Login(ctx, grpcClient)

	kafkaConsumer, err := messaging.NewKafkaConsumer("localhost:9092", fmt.Sprintf("my-group-%v", userID), *grpcClient, fmt.Sprintf("%v", userID))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Subscribe("test")

	for {
		var c Command

		fmt.Println("choose command")

		_, err = fmt.Scan(&c)
		if err != nil {
			log.Fatal("ой")
		}
		switch c {
		case WRITE:
			client.WriteMessage(ctx, grpcClient, userID)
		}
	}
}
