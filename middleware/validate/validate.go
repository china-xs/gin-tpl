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
					fmt.Println("err:", err.Error())
					i18nKey, params := getValidateKey(err.Error())
					fmt.Printf("i18nKey:%v,params:%v\n", i18nKey, params)
					msg := I18n.Tr(en, i18nKey, params)
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

func getValidateKey(str string) (string, map[string]interface{}) {
	var msgKey string
	params := make(map[string]interface{})
	if strings.Contains(str, "|") {
		for _, v := range strings.Split(str, "|") {
			if !strings.Contains(v, "invalid") {
				continue
			}
			child, cparams := getValidateKey(v)
			// 获取结果为空
			if child == "" {
				continue
			}
			if msgKey != "" {
				msgKey += "." + child
			} else {
				msgKey = child
			}
			for ck, cv := range cparams {
				params[ck] = cv
			}
		}
		return msgKey, params
	}
	key := str[strings.Index(str, "invalid ")+8:]
	if !strings.Contains(key, ":") {
		return msgKey, params
	}
	key = key[:strings.Index(key, ":")]
	if strings.Contains(key, "[") {
		params["key"] = key[strings.Index(key, "[")+1 : len(key)-1]
		key = key[:strings.Index(key, "[")]
	}
	cdn := getCondition(str)
	if cdn != nil {
		switch cdn.Key {
		case "between":
			tmp := str[strings.Index(str, cdn.Mst)+len(cdn.Mst):]
			r1 := strings.Contains(str, "runes")
			r2 := strings.Contains(str, "bytes")
			if r1 || r2 {
				strSlice := strings.Split(tmp, " ")
				if r1 {
					params["min_len"] = strSlice[0]
					params["max_len"] = strSlice[2]
				} else {
					params["min_bytes"] = strSlice[0]
					params["max_bytes"] = strSlice[2]
				}
			} else {
				var tk string
				if strings.Contains(str[strings.Index(str, cdn.Mst)+len(cdn.Mst):], "(") {
					tk = "gt"
				} else {
					tk = "gte"
				}
				params[tk] = tmp[1:strings.Index(tmp, ",")]
				if strings.Contains(str[strings.Index(str, cdn.Mst)+len(cdn.Mst):], ")") {
					tk = "lt"
				} else {
					tk = "lte"
				}
				params[tk] = tmp[strings.Index(tmp, ",")+2 : len(tmp)-1]
			}

		case "lt", "gt", "lte", "gte", "in", "not_in":
			tmp := str[strings.Index(str, cdn.Mst)+len(cdn.Mst):]
			params[cdn.Key] = tmp
		case "const":
			tmp := str[strings.Index(str, cdn.Mst)+len(cdn.Mst):]
			params["const"] = tmp
		case "len":
			tmp := str[strings.Index(str, cdn.Mst)+len(cdn.Mst):]
			strSlice := strings.Split(tmp, " ")
			if strings.Contains(tmp, "bytes") {
				cdn.Key = "len_bytes"
			} else {
				cdn.Key = "len"
			}
			params[cdn.Key] = strSlice[0]
		case "min_bytes", "max_bytes":
			tmp := str[strings.Index(str, cdn.Mst)+len(cdn.Mst):]
			strSlice := strings.Split(tmp, " ")
			if strSlice[1] == "runes" {
				if cdn.Key == "min_bytes" {
					cdn.Key = "min_len"
				} else {
					cdn.Key = "max_len"
				}
			}
			params[cdn.Key] = strSlice[0]
		case "repeated.between":
			tmp := str[strings.Index(str, cdn.Mst)+len(cdn.Mst):]
			strSlice := strings.Split(tmp, " ")
			params["min_items"] = strSlice[0]
			params["max_items"] = strSlice[2]
		}
		key += "." + cdn.Key
	}
	if msgKey != "" {
		msgKey += "." + key
	} else {
		msgKey = key
	}
	//fmt.Println("err:",str)
	//fmt.Println("params",params)
	return msgKey, params
}

func getCondition(str string) *msgkey {
	for _, v := range vts {
		if ok := strings.Contains(str, v.Mst); ok {
			return &v
		}
	}
	return nil
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
	{Mst: "value must be greater than or equal to ", Key: "gte"},
	{Mst: "value must be greater than ", Key: "gt"},
	{Mst: "value must be less than or equal to ", Key: "lte"},
	{Mst: "value must be less than ", Key: "lt"},
	{Mst: "value must be outside range ", Key: "between"}, // 没遇到
	{Mst: "value must be inside range ", Key: "between"},

	{Mst: "value length must be at least ", Key: "min_bytes"},
	{Mst: "value length must be at most ", Key: "max_bytes"},
	{Mst: "value length must be between ", Key: "between"},
	{Mst: "value length must be ", Key: "len"},
	{Mst: "value does not match regex pattern", Key: "pattern"},
	{Mst: "value does not have prefix", Key: "prefix"},
	{Mst: "value does not have suffix", Key: "suffix"},
	{Mst: "value does not contain substring", Key: "contains"},
	{Mst: "value contains substring", Key: "not_contains"},
	{Mst: "value must be a valid IP address", Key: "ip"},
	{Mst: "value must be a valid IPv4 address", Key: "ipv4"},
	{Mst: "value must be a valid IPv6 address", Key: "ipv6"},
	{Mst: "value must equal ", Key: "const"},
	{Mst: "value must be one of the defined enum values", Key: "enum"},
	{Mst: "value must be in list ", Key: "in"},
	{Mst: "value must not be in list ", Key: "not_in"},
	{Mst: "value is required", Key: "required"},
	// repeated 数组规则
	{Mst: "repeated value must contain unique items", Key: "repeated.unique"},
	{Mst: "value must contain exactly ", Key: "repeated.min_items"},      // 没遇到来
	{Mst: "value must contain at least", Key: "repeated.min_items"},      // 数组长度小于最小长度
	{Mst: "value must contain no more than ", Key: "repeated.max_items"}, //没遇到
	{Mst: "value must contain no more than ", Key: "repeated.max_items"}, // 数组长度超过最大只
	{Mst: "value must contain between ", Key: "repeated.between"},        // 区间

	{Mst: "value must be a valid email address", Key: "email"},
	{Mst: "value must be a valid hostname, or ip address", Key: "address"}, //优先级比 hostname 高
	{Mst: "value must be a valid hostname", Key: "hostname"},
	{Mst: "value must be absolute", Key: "uri"},
	{Mst: "value must be a valid URI", Key: "url"},
	{Mst: "value must be a valid UUID", Key: "uuid"},
}
