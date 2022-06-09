/**
 * @Author: EDZ
 * @Description:
 * @File: apiverifier
 * @Version: 1.0.0
 * @Date: 2022/6/9 11:37
 */

package apiverifier

import (
	gin_tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/china-xs/gin-tpl/pkg/api_sign"
	"github.com/gin-gonic/gin"
)

//md5 signer only
func ApiVerifier(options *api_sign.Options) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (reply interface{}, err error) {
			path := c.GetString(gin_tpl.OperationKey)
			var a *api_sign.ApiSign
			a = api_sign.NewApiSign(options)
			if !a.NeedCheck(path) {
				return handler(c, req)
			}
			err = a.Verifier(c.Request)
			if err != nil {
				return nil, err
			}

			return handler(c, req)
		}
	}
}
