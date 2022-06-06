// Package topsdk
// @author: xs
// @date: 2022/4/26
// @Description: 内部使用，非内部人员不用查看
package topsdk

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"
const LogKey = "qimen_sdk"

type (
	Options struct {
		AppKey    string `yaml:"appKey"` // 小程序appkey
		Secret    string `yaml:"secret"` // secret
		YjfKey    string `yaml:"yjfKey"`
		YjfSecret string `yaml:"yjfSecret"`
		// "http://qimen.api.taobao.com/router/qmtest"
		// "http://qimen.api.taobao.com/router/qm"
		Url          string `yaml:"url"`      // 访问域名
		SellerId     string `yaml:"sellerId"` // 店铺ID
		TargetAppKey string `yaml:"targetAppKey"`
		QmClient     TopClient
		Log          *zap.Logger
	}

	Qimen interface {
		OpenId2Ouid(ctx context.Context, openId string) (ouid string, err error)
		OpenId2OuidV0(ctx context.Context, openId string) (ouid string, err error)
	}
	OpenId2OuidReply struct {
		Code      string      `json:"code"`
		Data      string      `json:"data"`
		Success   bool        `json:"success"`
		Msg       string      `json:"msg"`
		Response  ResponseErr `json:"response"`
		RequestId string      `json:"reqeust_id"`
		//request_id
	}
	ResponseErr struct {
		Flag       string `json:"flag"`
		Code       int    `json:"code"`
		Message    string `json:"message"`
		SubMessage string `json:"sub_message"`
		RequestId  string `json:"request_id"`
	}
)

var ProviderSet = wire.NewSet(NewQimen)

func NewQimen(v *viper.Viper, log *zap.Logger) (Qimen, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("yjfQimen", o); err != nil {
		return nil, err
	}
	o.QmClient = NewDefaultTopClient(
		o.AppKey,
		o.Secret,
		o.Url,
		20000,
		20000)
	o.Log = log
	return o, err
}

var fileMap map[string]interface{}

func (qm Options) OpenId2Ouid(ctx context.Context, openId string) (ouid string, err error) {
	//jn := getToken(tt.args.yjfKey, tt.args.yjfSecret, "")
	jn := qm.getToken()
	paramMap := make(map[string]interface{})
	paramMap["bizValue"] = "crm-互动"
	paramMap["buyer"] = openId
	paramMap["description"] = jn
	paramMap["eventType"] = "1"
	paramMap["seller"] = qm.SellerId
	paramMap["target_app_key"] = qm.TargetAppKey

	jsonStr, err := qm.QmClient.Execute("qimen.taobao.miniapp.crm.connect", paramMap, fileMap)

	var reply OpenId2OuidReply
	l := log.WithCtx(ctx, qm.Log)
	if err != nil {
		l.Error(LogKey, zap.Error(err), zap.String("openId", openId))
		return "", err
	}
	// {"response":{"flag":"failure","code":41,"message":"Invalid arguments","sub_message":"buyer(openId)有误:AAHagWGKANKwnl6niApyldHh","request_id":"16lo5im0is8iv"}}
	// {"code":"0","data":"AAEzAT6jAAAN7QwAASEJK6lO","msg":"成功","success":true}
	json.Unmarshal([]byte(jsonStr), &reply)
	if reply.Response.Code > 0 || !reply.Success {
		var msg string
		msg = reply.Response.SubMessage
		if len(msg) == 0 {
			msg = reply.Response.Message
		}
		err = errors.New(msg)
		l.Error(LogKey, zap.Error(err), zap.String("reply", jsonStr))
		return
	}
	return reply.Data, nil
}

//
// OpenId2OuidV0
// @Description: 老版本链路 description=链路ID
// @receiver qm
// @param ctx
// @param openId
// @return ouid
// @return err
//
func (qm Options) OpenId2OuidV0(ctx context.Context, openId string) (ouid string, err error) {
	var traceId string
	if span := otelTrace.SpanContextFromContext(ctx); span.HasTraceID() {
		traceId = span.TraceID().String()
	} else {
		uid, _ := uuid.NewRandom()
		traceId = uid.String()
	}
	paramMap := make(map[string]interface{})
	paramMap["bizValue"] = "crm-互动"
	paramMap["buyer"] = openId
	paramMap["description"] = traceId
	paramMap["eventType"] = "1"
	paramMap["seller"] = qm.SellerId
	paramMap["target_app_key"] = qm.TargetAppKey
	jsonStr, err := qm.QmClient.Execute("qimen.taobao.miniapp.crm.connect", paramMap, fileMap)
	var reply OpenId2OuidReply
	l := log.WithCtx(ctx, qm.Log)
	if err != nil {
		l.Error(LogKey, zap.Error(err), zap.String("openId", openId))
		return "", err
	}
	// {"response":{"flag":"failure","code":41,"message":"Invalid arguments","sub_message":"buyer(openId)有误:AAHagWGKANKwnl6niApyldHh","request_id":"16lo5im0is8iv"}}
	// {"code":"0","data":"AAEzAT6jAAAN7QwAASEJK6lO","msg":"成功","success":true}
	json.Unmarshal([]byte(jsonStr), &reply)
	if reply.Response.Code > 0 || !reply.Success {
		var msg string
		msg = reply.Response.SubMessage
		if len(msg) == 0 {
			msg = reply.Response.Message
		}
		err = errors.New(msg)
		l.Error(LogKey, zap.Error(err), zap.String("reply", jsonStr))
		return
	}
	return reply.Data, nil

	return
}

func (o Options) getToken() (jn string) {
	at := time.Now().Format(TimeFormat)
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%v%v%v",
		o.YjfKey,
		o.YjfSecret,
		at)))
	token := hex.EncodeToString(h.Sum(nil))
	jn = fmt.Sprintf(`{"appkey":"%v","sign":"%v","timestamp":"%v"}`, o.YjfKey, token, at)
	return jn
}
