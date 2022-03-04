// Package gin_tpl
// @author: xs
// @date: 2022/3/4
// @Description: 数据解析
package gin_tpl

import "github.com/gin-gonic/gin"

type EncodeResponseFunc func(*gin.Context, interface{}, error)

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(c *gin.Context, obj interface{}, err error) {
	// 默认输出逻辑
	if err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, obj)
	}
}
