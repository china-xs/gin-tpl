// Package gin_tpl
// @author: xs
// @date: 2022/3/3
// @Description: gin_tpl,过滤器|gin中间件
package gin_tpl

import "github.com/gin-gonic/gin"

type FilterFunc func() gin.HandlerFunc
