package setup

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	kratos "github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/config/setting"
	"github.com/liangguifeng/kratos-app/internal/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var WatchConfigFields = make(map[string]*Field)

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

type Field struct {
	Name         string
	NacosKeyName string
	Value        reflect.Value
	Type         reflect.Type
}

// NewConfiger 获取nacos的client
func NewConfiger(app *kratos.Application) (*Configer, error) {
	err := checkDirOrCreate(app.Name)
	if err != nil {
		return nil, err
	}

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config.GetNacosNamespaceId(),
			TimeoutMs:           NACOS_TIMEOU_MS,
			NotLoadCacheAtStart: true,
			LogDir:              NACOS_LOG_DIR + "/" + app.Name,
			CacheDir:            NACOS_CACHE_DIR + "/" + app.Name,
			LogLevel:            NACOS_LOG_LEVEL,
		},
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(config.GetNacosAddress(), config.GetNacosEndpoint()),
		},
	})
	if err != nil {
		log.Panic(err)
	}
	return &Configer{client: client, currentNamespace: "public", DataId: app.Name}, nil
}

// CreateConfiger 创建配置
func (c *Configer) CreateConfiger(allKey string) error {
	ok, err := c.client.PublishConfig(vo.ConfigParam{
		DataId:  c.DataId,
		Group:   NACOS_DEFAULT_GROUP,
		Content: allKey,
	})
	if err != nil {
		return nil
	}
	if !ok {
		return errors.New("config create fail")
	}

	return nil
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

	// 如果获取到的key为空，可能是没有创建配置，此处创建默认配置
	if len(allKey) == 0 {
		allKey = `
{
  "Mysql.Host": "127.0.0.1",
  "Mysql.UserName": "root",
  "Mysql.Password": "root",
  "Mysql.DBName": "kratos-layout",
  "Mysql.Charset": "utf8mb4",
  "Mysql.MaxIdle": "",
  "Mysql.MaxOpen": "",
  "Mysql.Loc": "",
  "Mysql.MultiStatements": "",
  "Mysql.ConnMaxLifeSecond": "",
  "Mysql.ParseTime": "",
  "Mysql.Timeout": "",
  "Redis.Host": "127.0.0.1",
  "Redis.Password": "",
  "Redis.PoolNum": "",
  "Redis.DB": "0",
  "Redis.MaxIdle": ""
}
`
		strings.Replace(allKey, "kratos-layout", c.DataId, -1)
		err = c.CreateConfiger(allKey)
		if err != nil {
			log.Panic(err)
		}
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

// GetReflectFields 从v获取反射字段。
func GetReflectFields(section string, v interface{}) (map[string]*Field, error) {
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Type must be a struct")
	}

	result := make(map[string]*Field)
	fieldCnt := typeOf.NumField()
	for i := 0; i < fieldCnt; i++ {
		name := typeOf.Field(i).Name
		nacosKeyName := fmt.Sprintf("%s.%s", section, name)
		result[nacosKeyName] = &Field{
			Name:         name,
			NacosKeyName: nacosKeyName,
			Type:         typeOf.Field(i).Type,
			Value:        valueOf.Field(i),
		}
	}

	return result, nil
}

func SaveWatchConfigField(v interface{}, fields map[string]*Field) error {
	configValues := make(map[string]interface{})
	for nacosKeyName, field := range fields {
		var value interface{}
		switch field.Type.Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			value = kratos.Configer.GetIntValue(nacosKeyName, 0)
		case reflect.String:
			value = kratos.Configer.GetStringValue(nacosKeyName, "")
		case reflect.Bool:
			value = kratos.Configer.GetBoolValue(nacosKeyName, false)
		default:
			return fmt.Errorf("Current field type is not be supported")
		}
		configValues[field.Name] = value
	}
	if len(configValues) == 0 {
		return nil
	}

	if err := mapstructure.Decode(configValues, v); err != nil {
		return err
	}

	return nil
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

// WatchUpdateConfig 监听nacos配置
func (c *Configer) WatchUpdateConfig() {
	if len(WatchConfigFields) == 0 {
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				kratos.Logger.Errorf(context.Background(), "WatchUpdateConfig Recover err: %v", err)
			}
		}()

		err := c.client.ListenConfig(vo.ConfigParam{
			DataId: c.DataId,
			Group:  NACOS_DEFAULT_GROUP,
		})
		if err != nil {
			return
		}
	}()
}

// checkDirOrCreate 检查文件夹是否存在，不存在则创建
func checkDirOrCreate(serverName string) error {
	dir := NACOS_LOG_DIR + "/" + serverName
	exists := config.Exists(dir)
	if !exists {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dir = NACOS_CACHE_DIR + "/" + serverName
	exists = config.Exists(dir)
	if !exists {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
