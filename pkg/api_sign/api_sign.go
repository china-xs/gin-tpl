// Package api_sign
// @author: ekin
// @date: 2022/5/31
// @Description: api请求验签
package api_sign

import (
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/parkingwang/go-sign"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var ProviderSet = wire.NewSet(NewApiSign, NewOps)
var (
	ErrParseQuery   = errors.New("【api_sign】parse query error")
	ErrKeyMiss      = errors.New("【api_sign】some keys missing")
	ErrTimeout      = errors.New("【api_sign】sign timeout")
	ErrSignNotMatch = errors.New("【api_sign】sign not match")
)

type (
	Options struct {
		Secret    string   `yaml:"secret"`
		Path      []string `yaml:"path"`
		Prefix    []string `yaml:"prefix"`
		Whitelist []string `yaml:"whitelist"`
		Timeout   int64    `yaml:"timeout"`
	}
	ApiSign struct {
		secret        string
		prefixPath    []string
		matchPath     map[string]struct{}
		whitelistPath map[string]struct{}
		timeout       time.Duration
	}
)

func NewOps(v *viper.Viper) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	if err = v.UnmarshalKey("apisign", o); err != nil {
		return nil, err
	}
	return o, err
}

func NewApiSign(o *Options) *ApiSign {
	a := &ApiSign{
		prefixPath:    make([]string, 0),
		matchPath:     make(map[string]struct{}, 0),
		whitelistPath: make(map[string]struct{}, 0),
		timeout:       2 * time.Minute,
	}

	if o.Timeout > 0 {
		a.SetTimeout(time.Duration(o.Timeout) * time.Minute)
	}

	return a.SetSecret(o.Secret).
		SetPath(o.Path...).
		SetPrefix(o.Prefix...).
		SetWhitelist(o.Whitelist...)
}

//setting secret
func (a *ApiSign) SetSecret(secret string) *ApiSign {
	a.secret = secret
	return a
}

// setting timeout
func (a *ApiSign) SetTimeout(t time.Duration) *ApiSign {
	a.timeout = t
	return a
}

// setting prefix paths
func (a *ApiSign) SetPrefix(paths ...string) *ApiSign {
	a.prefixPath = append(a.prefixPath, paths...)
	return a
}

// setting whitelist paths
func (a *ApiSign) SetWhitelist(paths ...string) *ApiSign {
	for _, path := range paths {
		a.whitelistPath[path] = struct{}{}
	}
	return a
}

// setting match paths
func (a *ApiSign) SetPath(paths ...string) *ApiSign {
	for _, path := range paths {
		a.matchPath[path] = struct{}{}
	}
	return a
}

//need check sign
func (a *ApiSign) NeedCheck(path string) bool {

	//whitelist
	if _, exists := a.whitelistPath[path]; exists {
		return false
	}

	//matched path && prefix path
	hasPath := false
	if _, exists := a.matchPath[path]; exists {
		hasPath = true
	} else {
		for _, p := range a.prefixPath {
			if strings.HasPrefix(path, p) {
				hasPath = true
				break
			}
		}
	}
	return hasPath
}

//verifier the sign
func (a *ApiSign) Verifier(r *http.Request) error {
	uri := r.URL.RequestURI()
	verifier := sign.NewGoVerifier()
	if err := verifier.ParseQuery(uri); err != nil {
		fmt.Printf("SignVerifier parseQuery err:%v", err)
		return ErrParseQuery
	}

	//check timeout
	verifier.SetTimeout(a.timeout)
	if err := verifier.CheckTimeStamp(); nil != err {
		fmt.Printf("SignVerifier timeout err:%v", err)
		return ErrTimeout
	}

	//check sign
	signer := sign.NewGoSignerMd5()
	signer.SetBody(verifier.GetBodyWithoutSign())
	signer.SetAppSecretWrapBody(a.secret)
	sign := signer.GetSignature()
	if verifier.MustString("sign") != sign {
		fmt.Printf("SignVerifier sign not match source:%s sign:%s", signer.GetSignBodyString(), sign)
		return ErrSignNotMatch
	}

	return nil
}
