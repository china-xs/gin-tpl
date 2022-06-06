/**
 * @Author: ekin
 * @Description: jwt-token parser
 * @File: jwt_auth_test.go
 * @Version: 1.0.0
 * @Date: 2022/5/31 13:39
 */

package api_sign

import (
	"fmt"
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
		"mobile": "15014164031",
	})
	t.Log(s)
}

//验证失败回调函数
func TestSignVerifierFailed(t *testing.T) {
	handler := SignVerifier("sign-secret", WithUnsignedCallback(
		func(c *gin.Context, err error) {
			assert.NotNil(t, err)
			t.Log("callback err:", err)
			c.String(http.StatusBadRequest, "unsigned callback")
		}))

	router := gin.New()
	router.Use(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	router.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

//验证token成功
func TestSignVerifierSuccess(t *testing.T) {
	const key = "sign-secret"
	handler := SignVerifier(key,
		WithTimeout(10*time.Minute),
		WithMustHasFields("user_id"),
		WithUnsignedCallback(
			func(c *gin.Context, err error) {
				assert.NotNil(t, err)
				t.Log("callback err:", err)
				c.String(http.StatusBadRequest, "unsigned callback")
			}))

	router := gin.New()
	router.Use(handler)
	router.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	query := CreateSign(key, map[string]string{
		"user_id": "256",
	})
	t.Log("query", query)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/info?%s", query), nil)

	router.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

//验证确实必要字段
func TestSignVerifierMissKeys(t *testing.T) {
	const key = "sign-secret"
	handler := SignVerifier(key,
		WithTimeout(10*time.Minute),
		WithMustHasFields("user_id", "mobile"), //丢失键
		WithUnsignedCallback(
			func(c *gin.Context, err error) {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, ErrKeyMiss)
				c.String(http.StatusBadRequest, "miss key")
			}))

	router := gin.New()
	router.Use(handler)
	router.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	query := CreateSign(key, map[string]string{
		"user_id": "256",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/info?%s", query), nil)
	router.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "miss key", w.Body.String())
}

//验证超时
func TestSignVerifierTimeout(t *testing.T) {
	const key = "sign-secret"
	handler := SignVerifier(key,
		WithTimeout(1*time.Second),
		WithMustHasFields("user_id"), //丢失键
		WithUnsignedCallback(
			func(c *gin.Context, err error) {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, ErrTimeout)
				c.String(http.StatusBadRequest, "timeout")
			}))

	router := gin.New()
	router.Use(handler)
	router.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	query := CreateSign(key, map[string]string{
		"user_id": "256",
	})

	//模拟超时
	time.Sleep(3 * time.Second)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/info?%s", query), nil)
	router.ServeHTTP(w, req)
	t.Log(w.Code, w.Body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "timeout", w.Body.String())
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
