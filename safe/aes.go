package safe

import (
	"crypto/cipher"
	"crypto/aes"
	"bytes"
	"errors"
	"encoding/base64"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func AesDecrypt(dst, key string) (str string, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err == nil {
		blockMode := newECBDecrypter(block)
		crypted, err := base64.StdEncoding.DecodeString(dst)
		origData := make([]byte, len(crypted))
		if err == nil {
			blockMode.CryptBlocks(origData, crypted)
			origData = pkcs5UnPadding(origData)
			str = string(origData)
		}
	}
	return
}

func AesEncrypt(src, key string) (str string, err error) {
	if src == "" {
		err = errors.New("plain content empty")
	} else {
		block, err := aes.NewCipher([]byte(key))
		if err == nil {
			ecb := newECBEncrypter(block)
			content := []byte(src)
			content = pkcs5Padding(content, block.BlockSize())
			crypted := make([]byte, len(content))
			ecb.CryptBlocks(crypted, content)
			str = base64.StdEncoding.EncodeToString(crypted)
		}
	}
	return
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
