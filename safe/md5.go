package safe

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func GetMD5String(str string) string {
	data := []byte(str)
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
