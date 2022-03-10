package auth

import (
	"context"
	"fmt"
	pb "github.com/china-xs/gin-tpl/examples/blog/api/auth"
	"github.com/china-xs/gin-tpl/examples/blog/internal/data/dao/query"
	"github.com/go-redis/redis/v8"
	otelTrace "go.opentelemetry.io/otel/trace"
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
	//s.redis.Get()
	var trace string
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		trace = span.TraceID().String()
	}
	fmt.Printf("traceId:%v\n", trace)
	dept := query.Use(s.db).OaDepartments
	res, err := dept.WithContext(ctx).Where(dept.ID.Eq(1)).First()
	if err != nil {
		return nil, err
	}
	fmt.Printf("dept:%v\n", res)

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
