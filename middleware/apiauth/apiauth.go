/**
 * @Author: ekin
 * @Description:
 * @File: auth
 * @Version: 1.0.0
 * @Date: 2022/6/9 10:34
 */

package apiauth

import (
	gin_tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/china-xs/gin-tpl/pkg/jwt_auth"
	"github.com/gin-gonic/gin"
)

const (
	jwtAudience  = "aud"
	jwtExpire    = "exp"
	jwtId        = "jti"
	jwtIssueAt   = "iat"
	jwtIssuer    = "iss"
	jwtNotBefore = "nbf"
	jwtSubject   = "sub"
)

// Authorize is a jwt-token parser middleware.
func Authorize(options *jwt_auth.Options) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (reply interface{}, err error) {
			path := c.GetString(gin_tpl.OperationKey)
			var a *jwt_auth.JwtAuth
			a = jwt_auth.NewJwtAuth(options)
			if !a.NeedCheck(path) {
				return handler(c, req)
			}

			claims, err := a.Verifier(c.Request)
			if err != nil {
				return nil, err
			}

			for k, v := range claims {
				switch k {
				case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
				default:
					c.Set(k, v)
				}
			}

			return handler(c, req)
		}
	}
}
