// Package gin_tpl
// @author: xs
// @date: 2022/3/4
// @Description: 数据解析
package gin_tpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"net/http"
	"strings"
)

type EncodeResponseFunc func(*gin.Context, interface{}, error)

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(c *gin.Context, obj interface{}, err error) {
	// 默认输出逻辑
	if err != nil {
		if err.Error() == "EOF" {
			err = errors.BadRequest("VALIDATE", "body is null")
		}
		se := errors.FromError(err)
		var bufReply []byte
		bufReply, err = json.Marshal(se)
		w := c.Writer
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", strings.Join([]string{"application", "json"}, "/"))
		w.WriteHeader(int(se.Code))
		_, _ = w.Write(bufReply)
	} else {
		c.JSON(200, obj)
	}
}
