package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func Connect(config Config) (DB *gorm.DB, ErrorConnect error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database, config.Encoding)
	DB, ErrorConnect = gorm.Open(mysql.New(mysql.Config{
		DSN: dns, // data source name
		// DefaultStringSize:         256,                                                                        // default size for string fields
		// DisableDatetimePrecision:  true,                                                                       // disable datetime precision, which not supported before MySQL 5.6
		// DontSupportRenameIndex:    true,                                                                       // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		// DontSupportRenameColumn:   true,                                                                       // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		// SkipInitializeWithVersion: false,                                                                      // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Off log
		// Logger: logger.Default.LogMode(logger.Info), // On Notice for DEBUG
		// Logger: logger.Default.LogMode(logger.Warn), // On Notice log
	})
	if ErrorConnect != nil {
		panic(ErrorConnect)
	}

	sqlDB, _ := DB.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(1)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// sqlDB.SetMaxOpenConns(1500)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// sqlDB.Stats()

	return
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Encoding string
}
