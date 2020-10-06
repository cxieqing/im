package db

import (
	"fmt"
	"im/pkg/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB

func init() {
	config := config.NewConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MysqlUsername,
		config.MysqlPassword,
		config.MysqlHost,
		config.MysqlPort,
		config.MysqlDatabase,
	)
	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("Fatal error mysql connect: %s \n", err))
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, e := MysqlDB.DB()
	if e != nil {
		panic(fmt.Errorf("Fatal error get sql.DB : %s \n", e))
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}
