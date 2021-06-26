package dao

import (
	"log"
	"time"

	"export-server/pkg/conf"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis() error {
	c := conf.AppConf
	return ConnRedis(&redis.Options{
		Addr:        c.GetString("redis.addr"),
		DB:          c.GetInt("redis.db"),
		Password:    c.GetString("redis.pwd"),
		PoolSize:    c.GetInt("redis.pool_size"),
		MaxRetries:  c.GetInt("redis.max_reties"),
		IdleTimeout: c.GetDuration("redis.idle_timeout") * time.Millisecond,
	})
}

func ConnRedis(opt *redis.Options) error {
	Rdb = redis.NewClient(opt)
	_, err := Rdb.Ping(Rdb.Context()).Result()
	if err != nil {
		log.Printf("[dao] redis fail, err=%s", err)
		return err
	}
	return nil
}
