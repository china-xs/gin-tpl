// Package api_sign
// @author: ekin
// @date: 2022/5/31
// @Description: api请求验签
package api_sign

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/parkingwang/go-sign"
	"net/http"
	"time"
)

var (
	ErrParseQuery   = errors.New("parse query error")
	ErrKeyMiss      = errors.New("some keys missing")
	ErrTimeout      = errors.New("timeout")
	ErrSignNotMatch = errors.New("sign not match")
)

type (
	SignOptions struct {
		Callback      UnsignedCallback
		MustHasFields []string
		Timeout       time.Duration
	}

	// UnSignedCallback then callback
	UnsignedCallback func(c *gin.Context, err error)
	// SignOption defines the method to customize an SignOptions.
	SignOption func(opts *SignOptions)
)

//md5 signer only
func SignVerifier(secret string, opts ...SignOption) gin.HandlerFunc {
	var signOpts SignOptions
	for _, opt := range opts {
		opt(&signOpts)
	}

	return func(c *gin.Context) {
		verifier := sign.NewGoVerifier()
		if err := verifier.ParseQuery(c.Request.URL.RequestURI()); err != nil {
			unsigned(c, fmt.Errorf("%s:%w", err.Error(), ErrParseQuery), signOpts.Callback)
			return
		}

		//check needed fields
		if len(signOpts.MustHasFields) > 0 {
			if err := verifier.MustHasOtherKeys(signOpts.MustHasFields...); nil != err {
				unsigned(c, fmt.Errorf("%s:%w", err.Error(), ErrKeyMiss), signOpts.Callback)
				return
			}
		}

		//check timeout
		verifier.SetTimeout(signOpts.Timeout)
		if err := verifier.CheckTimeStamp(); nil != err {
			unsigned(c, fmt.Errorf("%s:%w", err.Error(), ErrTimeout), signOpts.Callback)
			return
		}

		//check sign
		signer := sign.NewGoSignerMd5()
		signer.SetBody(verifier.GetBodyWithoutSign())
		signer.SetAppSecretWrapBody(secret)
		sign := signer.GetSignature()
		if verifier.MustString("sign") != sign {
			fmt.Printf("sign string:%s sign:%s", signer.GetSignBodyString(), sign)
			unsigned(c, ErrSignNotMatch, signOpts.Callback)
			return
		}

		c.Next()
	}
}

// setting unsigned callback.
func WithUnsignedCallback(callback UnsignedCallback) SignOption {
	return func(opts *SignOptions) {
		opts.Callback = callback
	}
}

// setting timeout
func WithTimeout(t time.Duration) SignOption {
	return func(opts *SignOptions) {
		opts.Timeout = t
	}
}

// setting mustHasFields
func WithMustHasFields(fields ...string) SignOption {
	return func(opts *SignOptions) {
		opts.MustHasFields = fields
	}
}

//unauthorized process
func unsigned(c *gin.Context, err error, callback UnsignedCallback) {
	c.Abort()
	if callback != nil {
		callback(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusBadRequest)
}
