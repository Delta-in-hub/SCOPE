package utils

import (
	"os"
	"strconv"
)

// GetEnvOrDefault 获取环境变量，如果不存在则返回默认值
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvAsIntOrDefault 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func GetEnvAsIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
