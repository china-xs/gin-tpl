/**
 * @Author: ekin
 * @Description: jwt-token parser
 * @File: jwt_auth_test.go
 * @Version: 1.0.0
 * @Date: 2022/5/31 13:39
 */

package jwt_auth

import (
	tpl "github.com/china-xs/gin-tpl"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//验证失败回调函数
func TestPath(t *testing.T) {
	const path = "/match/info"
	const path1 = "/member/info1"
	const path2 = "/whitelist/info2"
	jwtAuth := NewJwtAuth()
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		jwtAuth.Path(path).
			Prefix("/member").
			Whitelist(path2).
			Authorize("test-secret"),
	)
	ops = append(ops,
		ms,                 // 中间件
		tpl.OpenApi(false), //在线文档
	)
	app := tpl.NewServer(ops...)

	app.Engine.GET(path, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return "ok info", nil
		})
		out, err := h(c, nil)
		app.Enc(c, out, err)
		return
	})

	app.Engine.GET(path1, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path1)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return "ok info1", nil
		})
		out, err := h(c, nil)
		app.Enc(c, out, err)
		return
	})

	app.Engine.GET(path2, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path2)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {
			return "ok info2", nil
		})
		out, err := h(c, nil)
		app.Enc(c, out, err)
		return
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	app.Engine.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", path1, nil)
	app.Engine.ServeHTTP(w1, req1)
	t.Log(w1.Code, w1.Body)
	assert.Equal(t, http.StatusInternalServerError, w1.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", path2, nil)
	app.Engine.ServeHTTP(w2, req2)
	t.Log(w2.Code, w2.Body)
	assert.Equal(t, http.StatusOK, w2.Code)
}

//验证失败
func TestAuthorizeFailed(t *testing.T) {
	const path = "/info"
	jwtAuth := NewJwtAuth()
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		jwtAuth.Path(path).Authorize("test-secret"),
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
	req, _ := http.NewRequest("GET", "/info", nil)
	app.Engine.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//验证token成功
func TestAuthorizeSuccess(t *testing.T) {
	const key = "test-secret"
	const path = "/info"
	jwtAuth := NewJwtAuth()
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		jwtAuth.Path(path).Authorize("test-secret"),
	)
	ops = append(ops,
		ms,                 // 中间件
		tpl.OpenApi(false), //在线文档
	)
	app := tpl.NewServer(ops...)

	app.Engine.GET(path, func(c *gin.Context) {
		c.Set(tpl.OperationKey, path)
		h := app.Middleware(func(c *gin.Context, req interface{}) (interface{}, error) {

			//验证token有效性及是否插入到gin上下文
			assert.Equal(t, "256", c.GetString("user_id"))
			return "ok", nil
		})
		out, err := h(c, nil)

		//不能报错
		assert.Nil(t, err)
		app.Enc(c, out, err)
		return
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	token, err := CreateToken(key, map[string]interface{}{
		"user_id": "256",
	}, 3600)
	t.Log("token", token)
	assert.Nil(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	app.Engine.ServeHTTP(w, req)
}

func CreateToken(secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
