package Common

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"server/Base"
)

type RedisTools struct{}

func (RedisTools) GetString(key string) string {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("Get", key))
	if err != nil {
		return ""
	}
	return res
}

func (RedisTools) SetString(key string, val string) string {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("Set", key, val, "EX", "60"))
	if err != nil {
		return ""
	}
	return res
}

func (RedisTools) SetInt(key string, val int) int {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	r, err := redis.Int(conn.Do("Set", key, val))
	if err != nil {
		fmt.Println("set  failed,", err)
		return -1
	}
	return r
}

func (RedisTools) GetInt(key string) int {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	r, err := redis.Int(conn.Do("Get", key))
	if err != nil {
		fmt.Println("get  failed,", err)
		return -1
	}
	return r
}
