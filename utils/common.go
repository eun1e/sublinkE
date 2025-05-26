package utils

import (
	"os"
	"strings"
)

// 检查环境
func CheckEnvironment() bool {
	APP_ENV := os.Getenv("APP_ENV")
	if APP_ENV == "" {
		// fmt.Println("APP_ENV环境变量未设置")
		return false
	}
	if strings.Contains(APP_ENV, "development") {
		// fmt.Println("你现在是开发环境")
		return true
	}
	// fmt.Println("你现在是生产环境")
	return false
}
