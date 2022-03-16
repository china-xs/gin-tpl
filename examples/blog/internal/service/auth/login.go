package auth

import (
	"context"
	"fmt"
	pb "github.com/china-xs/gin-tpl/examples/blog/api/auth"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/go-redis/redis/v8"
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
	//fmt.Printf("req:%+v\n", req)
	t := time.Now()
	//dept := query.Use(s.db).OaDepartments
	//dept.WithContext(ctx).Where(dept.ID.Eq(10)).Limit(1).Delete()

	db := s.db
	type Rest struct {
		Id            int32   `json:"id"`
		OrderAmount   int32   `json:"order_mounts" gorm:"column:order_mounts;not null"`
		OrderNum      float64 `json:"order_nums" gorm:"column:order_nums;not null"`
		PaymentAmount float64 `json:"payment_amounts" gorm:"column:payment_amounts;not null"`
	}
	var ss Rest

	db.Table("youshu_order").
		Where("date=?", "20220202").
		Select("id,CAST(sum(goods_amount_total) as DECIMAL(10, 3)) as order_mounts,sum(order_amount) as order_nums,sum(payment_amount) as payment_amounts").
		First(&ss)
	fmt.Printf("res:%+v\n", ss)
	//type Rest struct {
	//	MemberId int64 `json:"member_id"`
	//	TotalAge  int64`json:"total_age"`
	//	MpId int64 `json:"mp_id"`
	//}
	//var R Rest
	//db.WithContext(ctx).Table("t_mp_report").
	//	Where("member_id=?",10).
	//	Select("member_id,sum(`age`) as total_age,sum(mp_id) as mp_id").
	//	First(&R)
	//fmt.Printf("rest:%+v\n",R)

	////res, err := dept.WithContext(ctx).Where(dept.ID.Eq(5)).First()
	//if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	//	return nil, err
	//}
	lg := log.WithCtx(ctx, s.log)
	lg.Info("基础信息")
	lg.Error("错误信息")
	lg.Warn("警告信息")
	//fmt.Printf("dept:%v\n", res)
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
