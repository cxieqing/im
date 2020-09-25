package pkg

import (
	"fmt"
	"im/pkg/config"
	"sync"

	"github.com/garyburd/redigo/redis"
)

var instance *Redis
var once sync.Once

type Redis struct {
	conn redis.Conn
}

func GetInstance() *Redis {
	once.Do(func() {
		instance = new(Redis)
		config := config.NewConfig()
		host := config.RedisHost
		port := config.RedisPort
		conn, err := redis.Dial("tcp", host+":"+fmt.Sprint(port))
		if err != nil {
			panic(fmt.Errorf("Fatal error redis connect: %s \n", err))
		}
		instance.conn = conn
	})
	return instance
}

func NewRedis() *Redis {
	return GetInstance()
}
