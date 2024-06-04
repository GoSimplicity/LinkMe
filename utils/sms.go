package utils

import "math/rand"

// GenRandomCode 生成随机数
func GenRandomCode(length int) string {
	letters := []rune("123456789")
	code := make([]rune, length)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
