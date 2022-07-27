package utils

import "github.com/thinkeridea/go-extend/exunicode/exutf8"

// 字符串截取函数
func SubStrRuneIndexInString(s string, length int) string {
	n, _ := exutf8.RuneIndexInString(s, length)
	return s[:n]
}
