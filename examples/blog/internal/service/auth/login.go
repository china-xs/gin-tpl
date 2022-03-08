package auth

import (
	"context"
	pb "github.com/china-xs/gin-tpl/example/blog/api/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type LoginService struct {
	pb.UnimplementedLoginServer
}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (s *LoginService) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenReply, error) {
	t := time.Now()
	return &pb.GetTokenReply{
		Token:     "here with return a string token ",
		TokenType: "Bearer",
		ExpiresAt: timestamppb.New(t.Add(1 * time.Hour)),
	}, nil
}
