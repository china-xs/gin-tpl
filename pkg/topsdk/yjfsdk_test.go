// Package topsdk
// @author: xs
// @date: 2022/4/26
// @Description: topsdk
package topsdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"testing"
)

func getViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	// 任何需要将此配置添加到程序中的方法。
	var yamlExample = []byte(`
yjfQimen:
  appKey: ***
  secret: ****
  yjfKey: YJF_***
  yjfSecret: ***
  url: http://qimen.api.taobao.com/router/qmtest
  sellerId: **
  targetAppKey: ***
`)

	v.ReadConfig(bytes.NewBuffer(yamlExample))

	return v
}

func TestOptions_OpenId2Ouid(t *testing.T) {
	l := zap.NewExample()
	v := getViper()
	qm, err := NewQimen(v, l)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.TODO()
	ouid, err := qm.OpenId2Ouid(ctx, "AAEcjAK0ANdJu6_wxQxitG8i")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("ouid:", ouid)

}
