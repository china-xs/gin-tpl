// Package jwt_auth
// @author: ekin
// @date: 2022/5/31
// @Description: jwt token解析
package jwt_auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
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
	AuthorizeOptions struct {
		Callback UnauthorizedCallback
	}

	// UnauthorizedCallback then callback
	UnauthorizedCallback func(c *gin.Context, err error)
	// AuthorizeOption defines the method to customize an AuthorizeOptions.
	AuthorizeOption func(opts *AuthorizeOptions)
)

// Authorize is a jwt-token parser middleware.
func Authorize(secret string, opts ...AuthorizeOption) gin.HandlerFunc {
	var authOpts AuthorizeOptions
	for _, opt := range opts {
		opt(&authOpts)
	}

	return func(c *gin.Context) {
		//get token from http request
		var token *jwt.Token
		token, err := request.ParseFromRequest(
			c.Request,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			request.WithParser(newParser()))

		if err != nil {
			unauthorized(c, err, authOpts.Callback)
			return
		}

		if !token.Valid {
			unauthorized(c, ErrInvalidToken, authOpts.Callback)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized(c, ErrNoClaims, authOpts.Callback)
			return
		}

		for k, v := range claims {
			switch k {
			case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
			default:
				c.Set(k, v)
			}
		}

		c.Next()
	}
}

//jwt parser
func newParser() *jwt.Parser {
	return jwt.NewParser(jwt.WithJSONNumber())
}

// setting unauthorized callback.
func WithUnauthorizedCallback(callback UnauthorizedCallback) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.Callback = callback
	}
}

//unauthorized process
func unauthorized(c *gin.Context, err error, callback UnauthorizedCallback) {
	c.Abort()
	if callback != nil {
		callback(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusUnauthorized)
}
