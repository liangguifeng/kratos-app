package config

import (
	"os"
	"strconv"
)

const (
	// GO_ENV
	GO_ENV = "GO_ENV"
	// PROJECT_ENV
	PROJECT_ENV = "PROJECT_ENV"
	// PROJECT_LOGGER_PATH
	PROJECT_LOGGER_PATH = "PROJECT_LOGGER_PATH"
	// NACOS_ADDRESS
	NACOS_ADDRESS = "NACOS_ADDRESS"
	// NACOS_ENDPOINT
	NACOS_ENDPOINT = "NACOS_ENDPOINT"
	// NACOS_NAMESPACE_ID
	NACOS_NAMESPACE_ID = "NACOS_NAMESPACE_ID"
)

// GetBuildEnv 获取当前环境
func GetBuildEnv() string {
	value := os.Getenv(PROJECT_ENV)
	if value != "" {
		return value
	}

	value = os.Getenv(GO_ENV)
	if value != "" {
		return value
	}

	return ""
}

// GetNacosAddress
func GetNacosAddress() string {
	value := os.Getenv(NACOS_ADDRESS)
	if value != "" {
		return value
	}
	return ""
}

// GetNacosNamespaceId
func GetNacosNamespaceId() string {
	value := os.Getenv(NACOS_NAMESPACE_ID)
	if value != "" {
		return value
	}
	return ""
}

// GetNacosEndpoint
func GetNacosEndpoint() uint64 {
	value := os.Getenv(NACOS_ENDPOINT)
	if value != "" {
		intNum, _ := strconv.Atoi(value)
		int64Num := uint64(intNum)
		return int64Num
	}
	return 0
}

// GetProjectLoggerPath 获取日志根目录
func GetProjectLoggerPath() string {
	return os.Getenv(PROJECT_LOGGER_PATH)
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
