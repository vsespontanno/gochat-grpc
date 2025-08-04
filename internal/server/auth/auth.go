package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/vsespontanno/gochat-grpc/internal/models"
	"github.com/vsespontanno/gochat-grpc/internal/proto"
	pg "github.com/vsespontanno/gochat-grpc/internal/repository/pg"
)

type AuthService struct {
	proto.UnimplementedAuthServer
	pg *pg.UserStore
}

func NewAuthService(pg *pg.UserStore) *AuthService {
	return &AuthService{pg: pg}
}

func (s *AuthService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	var params models.User
	params.Email = req.GetEmail()
	params.Password = req.GetPassword()
	params.FirstName = req.GetFirstName()
	params.LastName = req.GetLastName()
	params.ID = uuid.New()
	err := s.pg.SaveUser(ctx, &params)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{UserId: int64(params.ID.ID())}, nil
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var authParams models.AuthParams
	authParams.Email = req.GetEmail()
	authParams.Password = req.GetPassword()
	user, err := s.pg.GetUserByEmail(ctx, authParams.Email)
	if err != nil {
		return nil, err
	}
	if user.Password != authParams.Password {
		return nil, err
	}
	return &proto.LoginResponse{Token: "poka vuyt ho"}, nil
}
