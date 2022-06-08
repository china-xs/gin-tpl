// Package jwt_auth
// @author: ekin
// @date: 2022/5/31
// @Description: jwt token解析
package jwt_auth

import (
	"errors"
	gin_tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"strings"
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

var (
	ErrInvalidToken = errors.New("invalid auth token")
	ErrNoClaims     = errors.New("no auth params")
)

type (
	JwtAuth struct {
		prefixPath    []string
		path          map[string]struct{}
		whitelistPath map[string]struct{}
	}
)

func NewJwtAuth() *JwtAuth {
	return &JwtAuth{
		prefixPath:    make([]string, 0),
		path:          make(map[string]struct{}, 0),
		whitelistPath: make(map[string]struct{}, 0),
	}
}

// Authorize is a jwt-token parser middleware.
func (a *JwtAuth) Authorize(secret string) middleware.Middleware {
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

			//get token from http request
			var token *jwt.Token
			token, err = request.ParseFromRequest(
				c.Request,
				request.AuthorizationHeaderExtractor,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				},
				request.WithParser(newParser()))

			if err != nil {
				return nil, err
			}

			if !token.Valid {
				return nil, ErrInvalidToken
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, ErrNoClaims
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

//jwt parser
func newParser() *jwt.Parser {
	return jwt.NewParser(jwt.WithJSONNumber())
}

// setting prefix paths
func (a *JwtAuth) Prefix(paths ...string) *JwtAuth {
	a.prefixPath = append(a.prefixPath, paths...)
	return a
}

// setting whitelist paths
func (a *JwtAuth) Whitelist(paths ...string) *JwtAuth {
	for _, path := range paths {
		a.whitelistPath[path] = struct{}{}
	}
	return a
}

// setting match paths
func (a *JwtAuth) Path(paths ...string) *JwtAuth {
	for _, path := range paths {
		a.path[path] = struct{}{}
	}
	return a
}
