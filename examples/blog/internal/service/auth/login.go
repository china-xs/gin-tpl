package auth

import (
	"context"
	"fmt"
	pb "github.com/china-xs/gin-tpl/examples/blog/api/auth"
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
	fmt.Println(req)
	return &pb.GetTokenReply{
		Token:     "here with return a string token ",
		TokenType: "Bearer",
		ExpiresAt: timestamppb.New(t.Add(1 * time.Hour)),
	}, nil
}

func (s *LoginService) GetInfo(ctx context.Context, in *pb.GetInfoRequest) (*pb.GetInfoReply, error) {
	fmt.Printf("in:%v\n", in)
	return &pb.GetInfoReply{}, nil
}
