package auth

import (
	"context"
	"fmt"
	pb "github.com/china-xs/gin-tpl/examples/blog/api/auth"
	"github.com/china-xs/gin-tpl/examples/blog/internal/data/dao/query"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type LoginService struct {
	pb.UnimplementedLoginServer
	log   *zap.Logger
	redis *redis.Client
	db    *gorm.DB
}

func NewLoginService(log *zap.Logger, db *gorm.DB, rdb *redis.Client) *LoginService {
	return &LoginService{
		log:   log,
		redis: rdb,
		db:    db,
	}
}

func (s *LoginService) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenReply, error) {
	t := time.Now()
	dept := query.Use(s.db).OaDepartments
	res, err := dept.WithContext(ctx).Where(dept.ID.Eq(5)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	lg := s.log.With(log.WithCtx(ctx)...)
	lg.Info("基础信息")
	lg.Error("错误信息")
	lg.Warn("警告信息")
	fmt.Printf("dept:%v\n", res)
	key := "gin-tpl:tmp"
	if t, err := s.redis.Exists(ctx, key).Result(); err != nil {
		fmt.Printf("redis:%v\n", t)
	}
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
