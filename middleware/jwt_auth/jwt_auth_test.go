/**
 * @Author: ekin
 * @Description: jwt-token parser
 * @File: jwt_auth_test.go
 * @Version: 1.0.0
 * @Date: 2022/5/31 13:39
 */

package jwt_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//验证失败回调函数
func TestAuthorizeFailed(t *testing.T) {
	handler := Authorize("test-secret", WithUnauthorizedCallback(
		func(c *gin.Context, err error) {
			assert.NotNil(t, err)
			c.String(http.StatusUnauthorized, "unauthorize callback")
		}))

	router := gin.New()
	router.Use(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	router.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

//验证token成功
func TestAuthorizeSuccess(t *testing.T) {
	const key = "test-secret"
	handler := Authorize(key, WithUnauthorizedCallback(
		func(c *gin.Context, err error) {
			assert.NotNil(t, err)
			c.String(http.StatusUnauthorized, "unauthorize callback")
		}))

	router := gin.New()

	//验证token有效性及是否插入到gin上下文
	router.Use(handler, func(c *gin.Context) {
		assert.Equal(t, "256", c.GetString("user_id"))
	})
	router.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	token, err := CreateToken(key, map[string]interface{}{
		"user_id": "256",
	}, 3600)
	t.Log("token", token)
	assert.Nil(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
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
