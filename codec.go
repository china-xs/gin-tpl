// Package gin_tpl
// @author: xs
// @date: 2022/3/4
// @Description: 数据解析
package gin_tpl

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*gin.Context, interface{}) error

type EncodeResponseFunc func(*gin.Context, interface{}, error)

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(c *gin.Context, v interface{}) error {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut:
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

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(c *gin.Context, obj interface{}, err error) {
	// 默认输出逻辑
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, obj)
	}
}
