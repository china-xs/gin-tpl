// Package validate
// @author: xs
// @date: 2022/3/9
// @Description: validate,描述
package validate

import (
	"fmt"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/kataras/i18n"
	"strings"
)

type validator interface {
	Validate() error
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error())
				}
			}
			return handler(c, req)
		}
	}
}

func Validator2I18n(I18n *i18n.I18n) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (interface{}, error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					en := c.Request.Header.Get("Accept-Language")
					if en == "" {
						en = "zh-CN"
					}
					res := strings.Split(err.Error(), " ")
					var i18nKey string
					i18nKey = strings.Replace(res[1], ":", "", 1)
					vstr := getValidate(err.Error())
					if vstr != "" {
						i18nKey += "." + vstr
						fmt.Println(i18nKey)
					}
					msg := I18n.Tr(en, i18nKey)
					if msg != "" {
						return nil, errors.BadRequest("VALIDATOR", msg)
					}
					return nil, errors.BadRequest("VALIDATOR", err.Error())
				}
			}
			return handler(c, req)
		}
	}
}

func getValidate(str string) string {
	for _, v := range vts {
		if ok := strings.Contains(str, v.Mst); ok {
			return v.Key
		}
	}
	return ""
}

type msgkey struct {
	Mst string
	Key string
}

// 新增 key 需要注意类型
var vts = []msgkey{
	{Mst: "value must be greater than or equal to", Key: "gte"},
	{Mst: "value must be greater than", Key: "gt"},
	{Mst: "value must be less than or equal to", Key: "lte"},
	{Mst: "value must be outside range", Key: "between"},
	{Mst: "value length must be at least", Key: "min_bytes"},
	{Mst: "value length must be at most", Key: "max_bytes"},
	{Mst: "value length must be between", Key: "between"},
	{Mst: "value length must be", Key: "len"},
	{Mst: "value does not match regex pattern", Key: "pattern"},
	{Mst: "value does not have prefix", Key: "prefix"},
	{Mst: "value does not have suffix", Key: "suffix"},
	{Mst: "value does not contain substring", Key: "contains"},
	{Mst: "value contains substring", Key: "not_contains"},
	{Mst: "value must be a valid IP address", Key: "ip"},
	{Mst: "value must be a valid IPv4 address", Key: "ipv4"},
	{Mst: "value must be a valid IPv6 address", Key: "ipv6"},
	{Mst: "value must equal", Key: "const"},
	{Mst: "value must be one of the defined enum values", Key: "enum"},
	{Mst: "value must be in list", Key: "in"},
	{Mst: "value must not be in list", Key: "not_in"},
	{Mst: "value is required", Key: "required"},
	{Mst: "value must contain exactly", Key: "min_items"},
	{Mst: "value must contain no more than", Key: "max_items"},
	{Mst: "repeated value must contain unique items", Key: "unique"},
	{Mst: "value must be a valid email address", Key: "email"},
	{Mst: "value must be a valid URI", Key: "url"},
}
