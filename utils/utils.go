package utils

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"runtime"
	"strings"
	"time"
)
import (
	"fmt"
	"math/rand"
	"strconv"
)

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	lstr := len(str) - 1
	for i := 0; i < l; i++ {
		n := GetRandomInt(0, lstr)
		result = append(result, bytes[n])
	}
	return string(result)
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomIntN(n int) int {
	return r.Intn(n)
}

func GetRandomInt(min, max int) int {
	sub := max - min + 1
	if sub <= 1 {
		return min
	}
	return min + GetRandomIntN(sub)
}

func GetPercentage(v, m, a int) float64 {
	x := float64(v) / float64(m) * 100
	f := fmt.Sprintf("%%.%df", a)
	t := fmt.Sprintf(f, x)
	x, err := strconv.ParseFloat(t, 10)
	if err != nil {
		return -1
	}
	return x
}

func GetPercentageString(v, m, a int, s bool) string {
	x := GetPercentage(v, m, a)
	f := fmt.Sprintf("%%.%df", a)
	if x < 0 {
		return ""
	}
	if s {
		return fmt.Sprintf(f, x) + "%"
	}
	return fmt.Sprintf(f, x)
}

func StringIndexOf(src, sub string, i int) int {
	if i < 0 || i > len(src)-1 {
		return -1
	}
	x := strings.Index(src[i:], sub) + i
	if x < i {
		return -1
	}
	return x
}

func UnixTimeBySeconds(s int64) time.Time {
	return time.Unix(s, 0)
}

func UnixTimeByMilliseconds(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}

func Utf8ToGbk(str string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), str)
	return result
}

func GbkToUtf8(str string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewDecoder(), str)
	return result
}

var newLine = func() string {
	s := "\n"
	switch runtime.GOOS {
	case "darwin":
		s = "\r"
		break
	case "windows":
		s = "\r\n"
		break
	}
	return s
}()

func NewLine() string {
	return newLine
}

func EnableAllCPU() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
