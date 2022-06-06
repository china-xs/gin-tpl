// Package api_sign
// @author: ekin
// @date: 2022/5/31
// @Description: api请求验签
package api_sign

import (
	"errors"
	"fmt"
	gin_tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"github.com/parkingwang/go-sign"
	"strings"
	"time"
)

var (
	ErrParseQuery   = errors.New("parse query error")
	ErrKeyMiss      = errors.New("some keys missing")
	ErrTimeout      = errors.New("timeout")
	ErrSignNotMatch = errors.New("sign not match")
)

type (
	ApiSign struct {
		prefixPath    []string
		path          map[string]struct{}
		whitelistPath map[string]struct{}
		timeout       time.Duration
	}
)

func NewApiSign() *ApiSign {
	return &ApiSign{
		prefixPath:    make([]string, 0),
		path:          make(map[string]struct{}, 0),
		whitelistPath: make(map[string]struct{}, 0),
		timeout:       2 * time.Minute,
	}
}

//md5 signer only
func (a *ApiSign) SignVerifier(secret string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (reply interface{}, err error) {
			path := c.GetString(gin_tpl.OperationKey)

			//whitelist
			if _, exists := a.whitelistPath[path]; exists {
				return handler(c, req)
			}

			//matched path && prefix path
			hasPath := false
			if _, exists := a.path[path]; exists {
				hasPath = true
			} else {
				for _, p := range a.prefixPath {
					if strings.HasPrefix(path, p) {
						hasPath = true
						break
					}
				}
			}

			if !hasPath {
				return handler(c, req)
			}

			verifier := sign.NewGoVerifier()
			if err := verifier.ParseQuery(c.Request.URL.RequestURI()); err != nil {
				fmt.Printf("SignVerifier parseQuery err:%v", err)
				return nil, ErrParseQuery
			}

			//check timeout
			verifier.SetTimeout(a.timeout)
			if err := verifier.CheckTimeStamp(); nil != err {
				fmt.Printf("SignVerifier timeout err:%v", err)
				return nil, ErrTimeout
			}

			//check sign
			signer := sign.NewGoSignerMd5()
			signer.SetBody(verifier.GetBodyWithoutSign())
			signer.SetAppSecretWrapBody(secret)
			sign := signer.GetSignature()
			if verifier.MustString("sign") != sign {
				fmt.Printf("SignVerifier sign not match source:%s sign:%s", signer.GetSignBodyString(), sign)
				return nil, ErrSignNotMatch
			}

			return handler(c, req)
		}
	}
}

// setting timeout
func (a *ApiSign) Timeout(t time.Duration) *ApiSign {
	a.timeout = t
	return a
}

// setting prefix paths
func (a *ApiSign) Prefix(paths ...string) *ApiSign {
	a.prefixPath = append(a.prefixPath, paths...)
	return a
}

// setting whitelist paths
func (a *ApiSign) Whitelist(paths ...string) *ApiSign {
	for _, path := range paths {
		a.whitelistPath[path] = struct{}{}
	}
	return a
}

// setting match paths
func (a *ApiSign) Path(paths ...string) *ApiSign {
	for _, path := range paths {
		a.path[path] = struct{}{}
	}
	return a
}
