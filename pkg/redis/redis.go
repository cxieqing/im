package redis

import (
	"fmt"
	"im/pkg/config"
	"sync"

	"github.com/go-redis/redis"
)

var instance *Redis
var once sync.Once

type Redis struct {
	Client *redis.Client
}

func GetInstance() *Redis {
	once.Do(func() {
		instance = new(Redis)
		config := config.NewConfig()

		redisdb := redis.NewClient(&redis.Options{
			Addr:     config.RedisHost + ":" + fmt.Sprint(config.RedisPort),
			Password: config.RedisAuth,
			DB:       0,
		})
		_, err := redisdb.Ping().Result()
		if err != nil {
			panic(fmt.Errorf("Fatal error redis connect: %s \n", err))
		}
		instance.Client = redisdb
	})
	return instance
}

func NewRedis() *Redis {
	return GetInstance()
}
