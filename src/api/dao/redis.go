package dao

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/garyburd/redigo/redis"
	"os"
	"path/filepath"
	// "sync"
	"time"
)

var redisPool *redis.Pool

// var redisOnce sync.Once

//init函数只会执行一次，线程安全
//Package initialization—variable initialization and the invocation of init functions—happens in a single goroutine, sequentially, one package at a time.
//n init function may launch other goroutines, which can run concurrently with the initialization code.
//However, initialization always sequences the init functions: it will not start the next init until the previous one has returned.
func init() {
	redisURL, err := getRedisURL()
	if err != nil {
		beego.Error("init redis config fail", err.Error())
		os.Exit(1)
	}
	redisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisURL)
			if err != nil {
				beego.Info("dial redis fail,", err.Error())
				return nil, err
			} else {
				beego.Info("dial redis success and return a connection")
			}
			return c, nil
		},
	}
}

//单例模式实现，lazy load
//但是调用的时候需:dao.GetRedisPool().Get()
// func GetRedisPool() *redis.Pool {
// 	redisOnce.Do(func() {
// 		redisURL, err := getRedisURL()
// 		if err == nil {
// 			redisPool = &redis.Pool{
// 				MaxIdle:     10,
// 				MaxActive:   100,
// 				IdleTimeout: 240 * time.Second,
// 				Dial: func() (redis.Conn, error) {
// 					c, err := redis.Dial("tcp", redisURL)
// 					if err != nil {
// 						beego.Error("dial redis fail,", err.Error())
// 						return nil, err
// 					} else {
// 						beego.Error("diao redis success and return a connection")
// 					}
// 					return c, nil
// 				},
// 			}
// 		}
// 	})
// 	return redisPool
// }

func getRedisURL() (string, error) {
	redisConfigPath := filepath.Join(beego.AppPath, "conf", "redis.ini")
	redisConfig, err := config.NewConfig("ini", redisConfigPath)
	if err != nil {
		beego.Error("new redis config error", err.Error())
		return "", err
	}

	envRunMode := os.Getenv("BEEGO_RUNMODE")
	if envRunMode == "" {
		envRunMode = beego.PROD
	}

	redisIP := redisConfig.String(envRunMode + "::ip")
	redisPort := redisConfig.String(envRunMode + "::port")
	redisAddr := redisIP + ":" + redisPort

	return redisAddr, nil
}

func SetCache(cacheName string, cacheVal interface{}) error {
	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", cacheName, cacheVal)

	if err != nil {
		beego.Error("set cache failed,", err.Error())
	}
	return err
}

//上次业务并不关心redis是否挂掉，如果为nil直接从db取
//error还是返回，由上层决定是否要处理
func GetCache(cacheName string) (interface{}, error) {
	conn := redisPool.Get()
	defer conn.Close()

	result, err := conn.Do("GET", cacheName)

	if err != nil {
		beego.Error("get cache failed,", err.Error())
		return nil, err
	}

	return result, err
}

func DelCache(cacheName string) error {
	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", cacheName)
	if err != nil {
		beego.Error("Del cache error,cacheName ", cacheName, " errmsg ", err.Error())
	}
	return err
}
