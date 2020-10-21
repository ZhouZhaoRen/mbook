package dynamicache

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool = nil
	MaxIdle   int         = 0
	MaxOpen   int         = 0
	ExpireSec int64       = 0
)

func InitCache() {
	addr := beego.AppConfig.String("dynamicache_addrstr")
	if len(addr) == 0 {
		addr = "127.0.0.1:6379"
	}
	if MaxIdle <= 0 {
		MaxIdle = 256
	}
	pass := beego.AppConfig.String("dynamicache_passwd")
	if len(pass) == 0 {
		pool = &redis.Pool{
			MaxIdle:     MaxIdle,
			MaxActive:   MaxOpen,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					addr,
					redis.DialReadTimeout(1*time.Second),
					redis.DialWriteTimeout(1*time.Second),
					redis.DialConnectTimeout(1*time.Second),
				)
			},
		}
	} else {
		pool = &redis.Pool{
			MaxIdle:     MaxIdle,
			MaxActive:   MaxOpen,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					addr,
					redis.DialReadTimeout(1*time.Second),
					redis.DialWriteTimeout(1*time.Second),
					redis.DialConnectTimeout(1*time.Second),
					redis.DialPassword(pass),
				)
			},
		}
	}
}

func rdsdo(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func WriteString(key string, value string) error {
	_, err := rdsdo("SET", key, value)
	beego.Debug("redis set:" + key + "-" + value)
	rdsdo("EXPIRE", key, ExpireSec)
	return err
}

func ReadString(key string) (string, error) {
	result, err := rdsdo("GET", key)
	beego.Debug("redis get:" + key)
	if nil == err {
		str, _ := redis.String(result, err)
		return str, nil
	} else {
		beego.Debug("redis get error:" + err.Error())
		return "", err
	}
}

 // 往Redis的string里面存的是一个结构体，先将结构体转为json数据再存入string
func WriteStruct(key string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if nil == err {
		return WriteString(key, string(data))
	} else {
		return nil
	}
}

// 先从string里面读出对应的数据，再通过json进行反序列化，再将这个转到obj里面
func ReadStruct(key string, obj interface{}) error {
	if data, err := ReadString(key); nil == err {
		return json.Unmarshal([]byte(data), obj)
	} else {
		return err
	}
}

func WriteList(key string, list interface{}, total int) error {
	realKeyList := key + "_list"
	realKeyCount := key + "_count"
	data, err := json.Marshal(list)
	if nil == err {
		WriteString(realKeyCount, strconv.Itoa(total))
		return WriteString(realKeyList, string(data))
	} else {
		return nil
	}
}

func ReadList(key string, list interface{}) (int, error) {
	realKeyList := key + "_list"
	realKeyCount := key + "_count"
	if data, err := ReadString(realKeyList); nil == err {
		totalStr, _ := ReadString(realKeyCount)
		total := 0
		if len(totalStr) > 0 {
			total, _ = strconv.Atoi(totalStr)
		}
		return total, json.Unmarshal([]byte(data), list)
	} else {
		return 0, err
	}
}
