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

const (
	_betweenRunes = "runes, inclusive"
	_betweenBytes = "bytes, inclusive"
	_betweenLen   = 16
	_strLen       = 6 // bytes |runes + 空格

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
					// validate 验证错误
					fmt.Println("err:", err.Error())
					i18nKey, params := getValidateKey(err.Error())
					// 转换对应 i18n key && 提供对应参数
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
	// 嵌套错误
	if strings.Contains(str, "|") {
		for _, v := range strings.Split(str, "|") {
			if getInvalid(v) == -1 {
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
	l := len(str)

	str = str[getInvalid(str)+1:]
	i := strings.IndexRune(str, ':')
	if i == -1 {
		return msgKey, params
	}
	key := str[:i]
	str = str[i+2:]
	l = len(str)
	// CreateUserRequest.Ids[4]
	if index := strings.IndexRune(key, '['); index != -1 {
		params["key"] = key[index+1 : len(key)-1]
		key = key[:index]
	}
	cdn := getCondition(str)
	if cdn != nil {
		switch cdn.Key {
		case "between":
			if t := str[l-_betweenLen:]; t == _betweenRunes || t == _betweenBytes {
				tmp := str[cdn.Len:]
				strSlice := strings.Split(tmp, " ")
				if t == _betweenRunes {
					params["min_len"] = strSlice[0]
					params["max_len"] = strSlice[2]
				} else {
					params["min_bytes"] = strSlice[0]
					params["max_bytes"] = strSlice[2]
				}
			} else {
				var tk string
				str = str[cdn.Len:]
				l = len(str)
				tk = "gt"
				if str[:1] == "[" {
					tk = "gte"
				}
				centerIndex := strings.IndexRune(str, ',')
				params[tk] = str[1:centerIndex]
				tk = "lt"
				if str[l-1:] != ")" {
					tk = "lte"
				}
				params[tk] = str[centerIndex+2 : l-1]
			}
		case "lt", "gt", "lte", "gte":
			params[cdn.Key] = str[cdn.Len:]
		case "in", "not_in":
			tmp := str[cdn.Len+1 : l-1]
			params[cdn.Key] = strings.Split(tmp, " ")
		case "const":
			params["const"] = str[cdn.Len:]
		case "len":
			if str[l-5:] == "bytes" {
				cdn.Key = "len_bytes"
			}
			params[cdn.Key] = str[cdn.Len : l-_strLen]
		case "min_bytes", "max_bytes":
			if str[l-5:] == "runes" {
				if cdn.Key == "min_bytes" {
					cdn.Key = "min_len"
				} else {
					cdn.Key = "max_len"
				}
			}
			s := str[:l-_strLen]
			params[cdn.Key] = s[strings.LastIndex(s, " ")+1:]
		case "repeated.between":
			t := str[cdn.Len:]
			strSlice := strings.Split(t, " ")
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
	return msgKey, params
}

func getInvalid(str string) int {
	if str[:8] == "invalid " {
		return 7
	} else if str[:19] == " caused by: invalid" {
		return 19
	}
	return -1
}

func getCondition(str string) *msgkey {
	for _, v := range vts {
		if ok := strings.Contains(str, v.Mst); ok {
			return &v
		}
	}
	return nil
}

type msgkey struct {
	Mst string
	Key string
	Len int
}

// 新增 key 需要注意类型
var vts = []msgkey{
	{Mst: "value must be greater than or equal to ", Key: "gte", Len: 39},
	{Mst: "value must be greater than ", Key: "gt", Len: 27},
	{Mst: "value must be less than or equal to ", Key: "lte", Len: 36},
	{Mst: "value must be less than ", Key: "lt", Len: 24},
	{Mst: "value must be outside range ", Key: "between", Len: 28}, // 没遇到
	{Mst: "value must be inside range ", Key: "between", Len: 27},

	{Mst: "value length must be at least ", Key: "min_bytes", Len: 30},
	{Mst: "value length must be at most ", Key: "max_bytes", Len: 29},
	{Mst: "value length must be between ", Key: "between", Len: 29},

	{Mst: "value length must be ", Key: "len", Len: 21}, // 禁止添加len_bytes

	{Mst: "value does not match regex pattern", Key: "pattern", Len: 34},
	{Mst: "value does not have prefix", Key: "prefix", Len: 26},
	{Mst: "value does not have suffix", Key: "suffix", Len: 26},
	{Mst: "value does not contain substring", Key: "contains", Len: 32},
	{Mst: "value contains substring", Key: "not_contains", Len: 24},
	{Mst: "value must be a valid IP address", Key: "ip", Len: 32},
	{Mst: "value must be a valid IPv4 address", Key: "ipv4", Len: 34},
	{Mst: "value must be a valid IPv6 address", Key: "ipv6", Len: 34},
	{Mst: "value must equal ", Key: "const", Len: 17},
	{Mst: "value must be one of the defined enum values", Key: "enum", Len: 44},
	{Mst: "value must be in list ", Key: "in", Len: 22},
	{Mst: "value must not be in list ", Key: "not_in", Len: 26},
	{Mst: "value is required", Key: "required", Len: 17},
	// repeated 数组规则
	{Mst: "repeated value must contain unique items", Key: "repeated.unique", Len: 40},
	{Mst: "value must contain exactly ", Key: "repeated.min_items", Len: 27},      // 没遇到来
	{Mst: "value must contain at least", Key: "repeated.min_items", Len: 27},      // 数组长度小于最小长度
	{Mst: "value must contain no more than ", Key: "repeated.max_items", Len: 32}, // 数组长度超过最大只
	{Mst: "value must contain between ", Key: "repeated.between", Len: 27},        // 区间

	{Mst: "value must be a valid email address", Key: "email", Len: 35},
	{Mst: "value must be a valid hostname, or ip address", Key: "address", Len: 45}, //优先级比 hostname 高
	{Mst: "value must be a valid hostname", Key: "hostname", Len: 30},
	{Mst: "value must be absolute", Key: "uri", Len: 22},
	{Mst: "value must be a valid URI", Key: "url", Len: 25},
	{Mst: "value must be a valid UUID", Key: "uuid", Len: 26},
}
