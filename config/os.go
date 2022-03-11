package config

import "os"

const (
	// ElasticStack APM Server URL
	ENV_APM_SERVER_URL = "ELASTIC_APM_SERVER_URL"
	// ElasticStack APM Server URLs
	ENV_APM_SERVER_URLS = "ELASTIC_APM_SERVER_URLS"
	// Apollo META Server URL
	ENV_APOLLO_META_SERVER_URL = "APOLLO_META_SERVER_URL"
	// ETCD V3 Server URL
	ENV_ETCDV3_SERVER_URL = "ETCDV3_SERVER_URL"
	// ETCD V3 Server URLs
	ENV_ETCDV3_SERVER_URLS = "ETCDV3_SERVER_URLS"
	// TENANT_NAMESPACE
	TENANT_NAMESPACE = "TENANT_NAMESPACE"
	// GO_ENV
	GO_ENV = "GO_ENV"
	// PROJECT_ENV
	PROJECT_ENV = "PROJECT_ENV"
	// PROJECT_LOGGER_PATH
	PROJECT_LOGGER_PATH = "PROJECT_LOGGER_PATH"
	// APOLLO_ACCESSKEY_SECRET
	APOLLO_ACCESSKEY_SECRET = "APOLLO_ACCESSKEY_SECRET"
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

// GetProjectLoggerPath 获取日志根目录
func GetProjectLoggerPath() string {
	return os.Getenv(PROJECT_LOGGER_PATH)
}
