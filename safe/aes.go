package safe

import (
	"crypto/cipher"
	"crypto/aes"
	"bytes"
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

func AesEncrypt(src, key []byte) (crypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err == nil {
		ecb := newECBEncrypter(block)
		content, err := pkcs5Padding(src, block.BlockSize())
		if err == nil {
			crypted = make([]byte, len(content))
			ecb.CryptBlocks(crypted, content)
		}
	}
	return
}

func AesDecrypt(dst, key []byte) (origData []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err == nil {
		blockMode := newECBDecrypter(block)
		origData = make([]byte, len(dst))
		if err == nil {
			blockMode.CryptBlocks(origData, dst)
			origData, err = pkcs5UnPadding(origData)
		}
	}
	return
}

func AesEncryptString(src, key string) (str string, err error) {
	data, err := AesEncrypt([]byte(src), []byte(key))
	if err == nil {
		str = base64.StdEncoding.EncodeToString(data)
	}
	return
}

func AesDecryptString(dst, key string) (str string, err error) {
	data, err := base64.StdEncoding.DecodeString(dst)
	if err == nil {
		data, err = AesDecrypt(data, []byte(key))
		if err == nil {
			str = string(data)
		}
	}
	return
}

func pkcs5Padding(ciphertext []byte, blockSize int) (data []byte, err error) {
	defer func() {
		t := recover()
		if t != nil {
			err = t.(error)
		}
	}()
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(ciphertext, padtext...)
	return
}

func pkcs5UnPadding(origData []byte) (data []byte, err error) {
	defer func() {
		t := recover()
		if t != nil {
			err = t.(error)
		}
	}()
	length := len(origData)
	unpadding := int(origData[length-1])
	data = origData[:length-unpadding]
	return
}
