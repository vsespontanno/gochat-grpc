package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vsespontanno/gochat-grpc/internal/client"
	"github.com/vsespontanno/gochat-grpc/internal/messaging"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
)

const endpoint = "localhost:8081"

type Command string

const (
	HELP  Command = "h"
	WRITE Command = "w"
)

func main() {

	ctx := context.Background()
	grpcClient, err := client.NewGRPCClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := messaging.NewKafkaConsumer("localhost:9092", "my-group", *grpcClient)
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
			var recipient, con string

			fmt.Println("enter recipient")
			_, err = fmt.Scan(&recipient)
			if err != nil {
				log.Fatal("ой 2")
			}

			fmt.Println("enter message")
			_, err = fmt.Scan(&con)
			if err != nil {
				log.Fatal("ой 3")
			}

			req := &proto.MessageRequest{
				Sender:    "me",
				Recipient: recipient,
				Content:   con,
			}
			_, err := grpcClient.SendMessage(ctx, req)
			if err != nil {
				fmt.Printf("failed to send message: %s", err.Error())
				continue
			}

			fmt.Println("message sent")

		}
	}
}
