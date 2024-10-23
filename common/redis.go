package common

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

var Redis *redis.Client

func initRedis() {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	poolsize := viper.GetInt("redis.poolsize")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 密码
		DB:       db,       // 数据库
		PoolSize: poolsize, // 连接池大小
	})

	rdb.Set(context.Background(), "1", 1145, time.Second*60*60)
	result, _ := rdb.Get(context.Background(), "1").Result()
	fmt.Println(result)
	Redis = rdb
}
