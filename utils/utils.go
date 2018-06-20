package utils

import "time"
import (
	"math/rand"
	"fmt"
	"strconv"
)

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
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
