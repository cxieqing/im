package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ImPort        int
	ImHost        string
	RpcPort       int
	RpcHost       string
	MaxClientNum  int
	ReadDeadline  int
	WriteDeadline int
	MaxGroupNum   int
	RedisHost     string
	RedisPort     int
	RedisAuth     string
	MysqlUsername string
	MysqlPassword string
	MysqlHost     string
	MysqlPort     string
	MysqlDatabase string
}

var instance *Config
var once sync.Once

func getInstance(path string) *Config {
	once.Do(func() {
		instance = new(Config)
		viper.SetConfigFile(path)
		viper.SetConfigType("json")
		//viper.AddConfigPath(".")
		err := viper.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		instance.ImPort = viper.GetInt("imPort")
		instance.ImHost = viper.GetString("imHost")
		instance.RpcPort = viper.GetInt("rpcPort")
		instance.RpcHost = viper.GetString("rpcHost")
		instance.MaxClientNum = viper.GetInt("maxClientNum")
		instance.ReadDeadline = viper.GetInt("readDeadline")
		instance.WriteDeadline = viper.GetInt("writeDeadline")
		instance.MaxGroupNum = viper.GetInt("maxGroupNum")
		instance.RedisHost = viper.GetString("redisHost")
		instance.RedisPort = viper.GetInt("redisPort")
		instance.RedisAuth = viper.GetString("redisAuth")
		instance.MysqlUsername = viper.GetString("mysqlUsername")
		instance.MysqlPassword = viper.GetString("mysqlPassword")
		instance.MysqlHost = viper.GetString("mysqlHost")
		instance.MysqlPort = viper.GetString("mysqlPort")
		instance.MysqlDatabase = viper.GetString("mysqlDatabase")
	})
	return instance
}

func NewConfig() *Config {
	return getInstance("D:/go/im/server.json")
}
