package setup

import (
	"encoding/json"
	kratos "github.com/liangguifeng/kratos-app"
	config2 "github.com/liangguifeng/kratos-app/config"
	"github.com/liangguifeng/kratos-app/config/setting"
	"github.com/liangguifeng/kratos-app/internal/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"strconv"
)

const (
	NACOS_DEFAULT_GROUP = "DEFAULT_GROUP"
	NACOS_TIMEOU_MS     = 5000
	NACOS_LOG_DIR       = "/tmp/nacos/log"
	NACOS_CACHE_DIR     = "/tmp/nacos/cache"
	NACOS_LOG_LEVEL     = "error"
)

type Configer struct {
	client           config_client.IConfigClient
	currentNamespace string
	DataId           string
}

// NewConfiger 获取nacos的client
func NewConfiger(app *kratos.Application) (*Configer, error) {
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config2.GetNacosNamespaceId(),
			TimeoutMs:           NACOS_TIMEOU_MS,
			NotLoadCacheAtStart: true,
			LogDir:              NACOS_LOG_DIR,
			CacheDir:            NACOS_CACHE_DIR,
			LogLevel:            NACOS_LOG_LEVEL,
		},
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(config2.GetNacosAddress(), config2.GetNacosEndpoint()),
		},
	})
	if err != nil {
		log.Panic(err)
	}
	return &Configer{client: client, currentNamespace: "public", DataId: app.Name}, nil
}

// GetAllKeys 获取所有nacos配置key
func (c *Configer) GetAllKeys() map[string]interface{} {
	allKey, err := c.client.GetConfig(vo.ConfigParam{
		DataId: c.DataId,
		Group:  NACOS_DEFAULT_GROUP,
	})
	if err != nil {
		log.Panic(err)
	}

	mapAllKey := make(map[string]interface{})
	err = json.Unmarshal([]byte(allKey), &mapAllKey)
	if err != nil {
		log.Panic(err)
	}

	return mapAllKey
}

// GetStringValue 获取string配置
func (c *Configer) GetStringValue(key, defaultValue string) string {
	mapAllKey := c.GetAllKeys()

	value := mapAllKey[key]
	if value == nil {
		return defaultValue
	}
	stringValue := value.(string)
	if stringValue == "" {
		return defaultValue
	}

	return stringValue
}

// GetStringValue 获取int配置
func (c *Configer) GetIntValue(key string, defaultValue int) int {
	value := c.GetStringValue(key, "")
	if value == "" {
		return defaultValue
	}

	vv, _ := strconv.Atoi(value)
	return vv
}

// GetBoolValue 获取bool配置
func (c *Configer) GetBoolValue(key string, defaultValue bool) bool {
	value := c.GetStringValue(key, "")
	if value == "" {
		return defaultValue
	}

	vv, _ := strconv.ParseBool(value)
	return vv
}

func (c *Configer) LoadAppConfig(app *kratos.Application) error {
	if setting.ServerSetting == nil {
		setting.ServerSetting = &setting.ServerSettingS{}
	}
	setting.MysqlSetting = &setting.MysqlSettingS{
		Host: kratos.Configer.GetStringValue(config.MYSQL_HOST, ""),
	}
	if setting.MysqlSetting.Host != "" {
		setting.MysqlSetting.UserName = kratos.Configer.GetStringValue(config.MYSQL_USERNAME, "")
		setting.MysqlSetting.Password = kratos.Configer.GetStringValue(config.MYSQL_PASSWORD, "")
		setting.MysqlSetting.DBName = kratos.Configer.GetStringValue(config.MYSQL_DBNAME, "")
		setting.MysqlSetting.Charset = kratos.Configer.GetStringValue(config.MYSQL_CHARSET, "")
		setting.MysqlSetting.Loc = kratos.Configer.GetStringValue(config.MYSQL_LOC, "")
		setting.MysqlSetting.MaxIdle = kratos.Configer.GetIntValue(config.MYSQL_MAXIDLE, 0)
		setting.MysqlSetting.MaxOpen = kratos.Configer.GetIntValue(config.MYSQL_MAXOPEN, 0)
		setting.MysqlSetting.ParseTime = kratos.Configer.GetBoolValue(config.MYSQL_PARSE_TIME, false)
		setting.MysqlSetting.Timeout = kratos.Configer.GetIntValue(config.MYSQL_TIMEOUT, 60)
		setting.MysqlSetting.MultiStatements = kratos.Configer.GetBoolValue(config.MYSQL_MULTI_STATEMENTS, false)
		setting.MysqlSetting.ConnMaxLifeSecond = kratos.Configer.GetIntValue(config.MYSQL_CONN_MAX_LIFE_SECOND, 0)
	}

	setting.RedisSetting = &setting.RedisSettingS{
		Host: kratos.Configer.GetStringValue(config.REDIS_HOST, ""),
	}
	if setting.RedisSetting.Host != "" {
		setting.RedisSetting.Password = kratos.Configer.GetStringValue(config.REDIS_PASSWORD, "")
		setting.RedisSetting.DB = kratos.Configer.GetIntValue(config.REDIS_DB, 0)
		setting.RedisSetting.MaxActive = kratos.Configer.GetIntValue(config.REDIS_MAXACTIVE, 0)
		setting.RedisSetting.MaxIdle = kratos.Configer.GetIntValue(config.REDIS_MAXIDLE, 0)
	}

	return nil
}
