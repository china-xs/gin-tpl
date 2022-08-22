// Package errors
// @author: xs
// @date: 2022/8/9
// @Description: errors
package errors

import (
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"testing"
)

func TestErrors(t *testing.T) {
	err := Newf(http.StatusBadRequest, "Crazy Thursday", "%s", "疯狂星期四")
	e1 := FromError(err)
	if !reflect.DeepEqual(e1.Code, int32(http.StatusBadRequest)) {
		t.Errorf("返回异常")
	}
	e := FromError(errors.New("test"))
	if !reflect.DeepEqual(e.Code, int32(UnknownCode)) {
		t.Errorf("no expect value: %v, but got: %v", e.Code, int32(UnknownCode))
	}
}
