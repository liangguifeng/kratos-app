package setting

var ServerSetting *ServerSettingS

type ServerSettingS struct {
	EndPoint             string
	IsRecordCallResponse bool
}

var MysqlSetting *MysqlSettingS

type MysqlSettingS struct {
	Host              string
	UserName          string
	Password          string
	DBName            string
	Charset           string
	MaxIdle           int
	MaxOpen           int
	Loc               string
	ConnMaxLifeSecond int
	MultiStatements   bool
	ParseTime         bool
	Timeout           int
}

var RedisSetting *RedisSettingS

type RedisSettingS struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	DB          int
}
