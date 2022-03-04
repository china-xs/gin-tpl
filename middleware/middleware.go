// Package middleware
// @author: xs
// @date: 2022/3/3
// @Description: middleware,中间件，核心处理-拦截、重置返回值
package middleware

import "github.com/gin-gonic/gin"

// Handler defines the handler invoked by Middleware.
type Handler func(c *gin.Context, req interface{}) (interface{}, error)

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(Handler) Handler
