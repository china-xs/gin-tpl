// Package crypt2go test
// @author: xs
// @date: 2022/7/19
// @Description: crypt2go
package crypt2go

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	key := "2R[<)NcD)^H4GDv."
	mobile := "13429030111"
	ct, err := AesEncrypt([]byte(mobile), []byte(key))
	if err != nil {
		t.Error(err)
		return
	}
	// us0fcdmiP7ZQVHJPx6Lhng==
	// us0fcdmiP7ZQVHJPx6Lhng==
	fmt.Println("base64:", base64.StdEncoding.EncodeToString(ct))
	buf, err := AesDecrypt(ct, []byte(key))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Recovered plaintext: %s\n", buf)
	return

}
