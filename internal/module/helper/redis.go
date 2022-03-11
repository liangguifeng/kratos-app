package helper

import (
	"fmt"
	"github.com/liangguifeng/kratos-app/config/setting"
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewRedisConn(redisSetting *setting.RedisSettingS) (*redisConn, error) {
	if redisSetting == nil {
		return nil, fmt.Errorf("redisSetting is nil")
	}
	if redisSetting.Host == "" {
		return nil, fmt.Errorf("lack of redisSetting.Host")
	}
	if redisSetting.Password == "" {
		return nil, fmt.Errorf("lack of redisSetting.Password")
	}

	maxIdle := 10
	maxActive := 15
	idleTimeout := 240
	if redisSetting.MaxActive > 0 && redisSetting.MaxIdle > 0 {
		maxIdle = redisSetting.MaxIdle
		maxActive = redisSetting.MaxActive
	}
	if redisSetting.IdleTimeout > 0 {
		idleTimeout = redisSetting.IdleTimeout
	}
	redisPool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisSetting.Host)
			if err != nil {
				return nil, err
			}
			if redisSetting.Password != "" {
				if _, err := c.Do("AUTH", redisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if redisSetting.DB > 0 {
				if _, err := c.Do("SELECT", redisSetting.DB); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &redisConn{client: redisPool}, nil
}

type redisConn struct {
	client *redis.Pool
}

func (r *redisConn) GetClient() *redis.Pool {
	return r.client
}

func (r *redisConn) ActiveCount() int {
	return r.client.ActiveCount()
}

func (r *redisConn) Close() error {
	return r.client.Close()
}

func (r *redisConn) IdleCount() int {
	return r.client.IdleCount()
}

func (r *redisConn) Stats() redis.PoolStats {
	return r.client.Stats()
}
