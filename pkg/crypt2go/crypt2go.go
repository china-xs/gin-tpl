// Package crypt2go 对称加密常规版本 与mysql AES_ENCRYPT('', 'key') 保持一致
// @author: xs
// @date: 2022/7/19
// @Description: crypt2go
package crypt2go

import (
	"crypto/aes"
	"github.com/china-xs/gin-tpl/pkg/crypt2go/ecb"
	"github.com/china-xs/gin-tpl/pkg/crypt2go/padding"
)

//
// AesEncrypt aes-加密
// @param plaintext
// @param key
// @return []byte
// @return error
//
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	plaintext, err = padder.Pad(plaintext) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		return nil, err
	}
	ct := make([]byte, len(plaintext))
	mode.CryptBlocks(ct, plaintext)
	return ct, nil
}

//
// AesDecrypt aes-解密
// @param plaintext
// @param key
// @return []byte
// @return error
//
func AesDecrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(plaintext))
	mode.CryptBlocks(pt, plaintext)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	return padder.Unpad(pt) // unpad plaintext after decryption
}
