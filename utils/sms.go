package utils

import (
	"math/rand"
	"regexp"
)

// GenRandomCode 生成随机数
func GenRandomCode(length int) string {
	letters := []rune("123456789")
	code := make([]rune, length)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

// IsValidNumber 检查给定的字符串是否符合中国的手机号格式
func IsValidNumber(number string) bool {
	// 中国的手机号通常是以1开头的11位数字
	// 这个正则表达式匹配以1开头，第二位是3、4、5、6、7、8、9中的一个，后面跟着9位数字的字符串
	pattern := `^1[3456789]\d{9}$`
	matched, _ := regexp.MatchString(pattern, number)
	return matched
}
