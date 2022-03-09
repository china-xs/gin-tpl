// Package service
// @author: xs
// @date: 2022/3/7
// @Description: service,描述
package service

import (
	srvAuth "github.com/china-xs/gin-tpl/examples/blog/internal/service/auth"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	srvAuth.NewLoginService,
)
