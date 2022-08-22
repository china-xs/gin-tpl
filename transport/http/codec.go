// Package http
// @author: xs
// @date: 2022/8/5
// @Description: http
package http

import (
	"github.com/china-xs/gin-tpl/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*gin.Context, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(*gin.Context, interface{}, error)

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(c *gin.Context, v interface{}) error {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch: // default Content-Type:json
		if err := c.ShouldBindBodyWith(v, binding.JSON); err != nil {
			return errors.New(http.StatusBadRequest, "bindBody", err.Error())
		}
	case http.MethodGet, http.MethodDelete:
		if err := c.ShouldBindQuery(v); err != nil {
			return errors.New(http.StatusBadRequest, "bindQuery", err.Error())
		}
	}
	return nil
}

// AnyRequestDecoder decodes the request body to object.
func AnyRequestDecoder(c *gin.Context, v interface{}) error {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		bin := getBindingBody(c)
		if err := c.ShouldBindBodyWith(v, bin); err != nil {
			return err
		}
	case http.MethodGet, http.MethodDelete:
		if err := c.ShouldBindQuery(v); err != nil {
			return err
		}
	}
	return nil
}

// getBindingBody 获取绑定类型
func getBindingBody(c *gin.Context) binding.BindingBody {
	b := binding.Default(c.Request.Method, c.ContentType())
	var bin binding.BindingBody
	switch b.Name() {
	case "json":
		bin = binding.JSON
	case "xml":
		bin = binding.XML
	case "yaml":
		bin = binding.YAML
	case "protobuf":
		bin = binding.ProtoBuf
	case "msgpack":
		bin = binding.MsgPack
	default:
		bin = binding.JSON
	}
	return bin
}

type Resp struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Reason string      `json:"reason"`
	Data   interface{} `json:"data"`
}

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(c *gin.Context, obj interface{}, err error) {
	var resp Resp
	resp.Msg = "请求成功"
	resp.Reason = "success"
	if err != nil {
		er1 := errors.FromError(err)
		resp.Code = int(er1.Code)
		resp.Msg = er1.Message
		resp.Reason = er1.Reason
		c.JSON(int(er1.Code), resp)
		return
	}
	resp.Data = obj
	c.JSON(http.StatusOK, resp)
	return
}
