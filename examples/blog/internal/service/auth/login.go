package auth

import (
	"context"
	"fmt"
	pb "github.com/china-xs/gin-tpl/examples/blog/api/auth"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type LoginService struct {
	pb.UnimplementedLoginServer
	log *zap.Logger
}

func NewLoginService(log *zap.Logger) *LoginService {
	return &LoginService{
		log: log,
	}
}

func (s *LoginService) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenReply, error) {
	t := time.Now()
	fmt.Println(req)
	s.log.Info("msg")
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
