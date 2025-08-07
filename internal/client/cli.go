package client

import (
	"context"
	"fmt"
	"log"

	"github.com/vsespontanno/gochat-grpc/internal/models"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
)

func Register(ctx context.Context, grpcClient *GRPCClient) error {
	var email, password, firstName, lastName string
	fmt.Println("enter email")
	_, err := fmt.Scan(&email)
	if err != nil {
		log.Fatal("ой 2")
	}

	fmt.Println("enter password")
	_, err = fmt.Scan(&password)
	if err != nil {
		log.Fatal("ой 3")
	}

	fmt.Println("enter first name")
	_, err = fmt.Scan(&firstName)
	if err != nil {
		log.Fatal("ой 4")
	}

	fmt.Println("enter last name")
	_, err = fmt.Scan(&lastName)
	if err != nil {
		log.Fatal("ой 5")
	}
	params := models.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	errors, err := params.Validate()
	if err != nil {
		for _, err := range errors {
			fmt.Println(err)
		}
		fmt.Println("Your registration failed. Please try again.")
		return err
	}
	user, err := models.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}

	req := &proto.RegisterRequest{
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	resp, err := grpcClient.Register(ctx, req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Registration successful. User ID:", resp.UserId)
	fmt.Println("Now you need to log in with your credentials: ")
	return nil
}

func Login(ctx context.Context, grpcClient *GRPCClient) int64 {
	var email, password string
	fmt.Println("enter email")
	_, err := fmt.Scan(&email)
	if err != nil {
		log.Fatal("ой 2")
	}
	fmt.Println("enter password")
	_, err = fmt.Scan(&password)
	if err != nil {
		log.Fatal("ой 3")
	}
	resp, err := grpcClient.Login(ctx, &proto.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Login successful. User ID:", resp.UserId)
	return resp.UserId
}

func WriteMessage(ctx context.Context, grpcClient *GRPCClient, sender int64) {
	var recipient, con string
	fmt.Println("enter recipient")

	_, err := fmt.Scan(&recipient)
	if err != nil {
		log.Fatal("ой 2")
	}
	fmt.Println("enter message")
	_, err = fmt.Scan(&con)
	if err != nil {
		log.Fatal("ой 3")
	}

	req := &proto.MessageRequest{
		Sender:    fmt.Sprintf("%v", sender),
		Recipient: recipient,
		Content:   con,
	}
	_, err = grpcClient.SendMessage(ctx, req)
	if err != nil {
		fmt.Printf("failed to send message: %s", err.Error())
	}

	fmt.Println("message sent")

}
