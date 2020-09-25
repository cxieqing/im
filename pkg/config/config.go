package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Port          int
	Host          string
	MaxClinetNum  int
	ReadDeadline  int64
	WriteDeadline int64
	MaxGroupNum   int
	RedisHost     string
	RedisPort     int
	RedisAuth     string
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
		instance.Port = viper.GetInt("Port")
		instance.Host = viper.GetString("Host")
		instance.MaxClinetNum = viper.GetInt("MaxClinetNum")
		instance.ReadDeadline = int64(viper.GetInt("ReadDeadline"))
		instance.WriteDeadline = int64(viper.GetInt("WriteDeadline"))
		instance.MaxGroupNum = viper.GetInt("MaxGroupNum")
		instance.RedisHost = viper.GetString("RedisHost")
		instance.RedisPort = viper.GetInt("RedisPort")
		instance.RedisAuth = viper.GetString("RedisAuth")
	})
	return instance
}

func NewConfig() *Config {
	return getInstance("./server.json")
}
