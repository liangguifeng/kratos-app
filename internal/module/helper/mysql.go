package helper

import (
	"database/sql"
	"fmt"
	"github.com/liangguifeng/kratos-app/config/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"time"
)

func NewMySQLConn(mysqlSetting *setting.MysqlSettingS) (*mysqlConn, error) {
	if mysqlSetting == nil {
		return nil, fmt.Errorf("mysqlSetting is nil")
	}
	if mysqlSetting.UserName == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.UserName")
	}
	if mysqlSetting.Password == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Password")
	}
	if mysqlSetting.Host == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Host")
	}
	if mysqlSetting.DBName == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.DBName")
	}
	if mysqlSetting.Charset == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Charset")
	}
	if mysqlSetting.Loc == "" {
		mysqlSetting.Loc = "Local"
	} else {
		mysqlSetting.Loc = url.QueryEscape(mysqlSetting.Loc)
	}

	sqlDB, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=%s&parseTime=%t&loc=%s&multiStatements=%t&timeout=%ds",
		mysqlSetting.UserName,
		mysqlSetting.Password,
		mysqlSetting.Host,
		mysqlSetting.DBName,
		mysqlSetting.Charset,
		mysqlSetting.ParseTime,
		mysqlSetting.Loc,
		mysqlSetting.MultiStatements,
		mysqlSetting.Timeout,
	))
	if err != nil {
		return nil, err
	}

	maxIdle := 10
	maxOpen := 30
	if mysqlSetting.MaxOpen > 0 && mysqlSetting.MaxIdle > 0 {
		maxIdle = mysqlSetting.MaxIdle
		maxOpen = mysqlSetting.MaxOpen
	}
	if mysqlSetting.ConnMaxLifeSecond > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(mysqlSetting.ConnMaxLifeSecond) * time.Second)
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &mysqlConn{client: db}, nil
}

type mysqlConn struct {
	client *gorm.DB
}

func (r *mysqlConn) GetClient() *gorm.DB {
	return r.client
}
