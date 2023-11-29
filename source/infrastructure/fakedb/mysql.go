package fakedb

import (
	"gorm.io/gorm"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
)

// NewMysql : tạo một client mysql db test để viết mock test cho repo
func NewMysql() (DB *gorm.DB, ErrorConnect error) {
	cfg := config.NewConfig()
	return mysql.Connect(mysql.Config{
		Username: cfg.Mysql.Username,
		Password: cfg.Mysql.Password,
		Host:     "127.0.0.1",
		Port:     "3306",
		// Database: cfg.Mysql.Database + "_fakedb",
		Database: cfg.Mysql.Database,
		Encoding: cfg.Mysql.Encoding,
	})
}

func NewMysqlAff() (DB *gorm.DB, ErrorConnect error) {
	cfg := config.NewConfig()
	return mysql.Connect(mysql.Config{
		Username: cfg.Mysql.Username,
		Password: cfg.Mysql.Password,
		Host:     cfg.Mysql.Host,
		Port:     cfg.Mysql.Port,
		// Host:     "127.0.0.1",
		// Port:     "3306",
		// Database: cfg.Mysql.Database + "_fakedb",
		Database: "apac_aff",
		Encoding: cfg.Mysql.Encoding,
	})
}

// NewCache : tạo một client caching aerospike db test để viết mock test cho repo
func NewCache() (client caching.Cache, err error) {
	return caching.NewAerospike(caching.Config{
		Host:      "127.0.0.1",
		Port:      3000,
		Namespace: "disk_cache",
	})
}
