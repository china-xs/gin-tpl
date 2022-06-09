/**
 * @Author: ekin
 * @Description: jwt-token parser
 * @File: jwt_auth_test.go
 * @Version: 1.0.0
 * @Date: 2022/5/31 13:39
 */

package apiverifier

import (
	"fmt"
	tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/pkg/api_sign"
	"github.com/gin-gonic/gin"
	"github.com/parkingwang/go-sign"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	appKey = "9d8a121ce581499d"
)

//新建sign
func TestCreateSign(t *testing.T) {
	s := CreateSign("secret", map[string]string{
		"mobile": "15014164032",
	})
	t.Log(s)
}

//验证签名成功
func TestSignVerifierSuccess(t *testing.T) {
	const secret = "test-secret"
	const path = "/info"
	options := &api_sign.Options{
		Secret:    secret,
		Path:      []string{path},
		Prefix:    []string{"/member"},
		Whitelist: []string{},
		Timeout:   2,
	}
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		ApiVerifier(options),
	)
	ops = append(ops,
		ms,                 // 中间件
		tpl.OpenApi(false), //在线文档
	)
	app := tpl.NewServer(ops...)

	app.Engine.GET(path, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		})
		out, err := h(c, nil)
		app.Enc(c, out, err)
		return
	})

	w := httptest.NewRecorder()
	query := CreateSign(secret, map[string]string{
		"user_id": "256",
	})
	t.Log("query", query)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/info?%s", query), nil)
	app.Engine.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
}

//验证超时
func TestSignVerifierTimeout(t *testing.T) {
	const secret = "test-secret"
	const path = "/info"
	options := &api_sign.Options{
		Secret:    secret,
		Path:      []string{path},
		Prefix:    []string{"/member"},
		Whitelist: []string{},
		Timeout:   2,
	}
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		ApiVerifier(options),
	)
	ops = append(ops,
		ms,                 // 中间件
		tpl.OpenApi(false), //在线文档
	)
	app := tpl.NewServer(ops...)

	app.Engine.GET(path, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		})
		out, err := h(c, nil)
		assert.EqualError(t, api_sign.ErrTimeout, "timeout")
		app.Enc(c, out, err)
		return
	})

	w := httptest.NewRecorder()
	query := CreateSign(secret, map[string]string{
		"user_id": "256",
	})
	t.Log("query", query)

	//模拟三秒后访问
	time.Sleep(3 * time.Second)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/info?%s", query), nil)
	app.Engine.ServeHTTP(w, req)
}

//验证签名不匹配
func TestSignVerifierNotMatch(t *testing.T) {
	const secret = "test-secret"
	const path = "/info"
	options := &api_sign.Options{
		Secret:    secret,
		Path:      []string{path},
		Prefix:    []string{"/member"},
		Whitelist: []string{},
		Timeout:   2,
	}
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		ApiVerifier(options),
	)
	ops = append(ops,
		ms,                 // 中间件
		tpl.OpenApi(false), //在线文档
	)
	app := tpl.NewServer(ops...)

	app.Engine.GET(path, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		})
		out, err := h(c, nil)
		assert.EqualError(t, api_sign.ErrSignNotMatch, "sign not match")
		app.Enc(c, out, err)
		return
	})

	w := httptest.NewRecorder()
	query := CreateSign(secret, map[string]string{
		"user_id": "256",
	})
	t.Log("query", query)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/info?a=1&%s", query), nil)
	app.Engine.ServeHTTP(w, req)
}

func CreateSign(secretKey string, payloads map[string]string) string {
	signer := sign.NewGoSignerMd5()
	signer.SetAppId(appKey)
	signer.SetTimeStamp(time.Now().Unix())
	signer.SetNonceStr(signer.GetNonceStr())
	for k, v := range payloads {
		signer.AddBody(k, v)
	}
	signer.SetAppSecretWrapBody(secretKey)
	return signer.GetSignedQuery()
}
