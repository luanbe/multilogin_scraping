package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewConnectionDB(driverDB string, database string, host string, user string, password string, port int) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	switch viper.GetString("database.logger_level") {
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "silent":
		logLevel = logger.Silent
	default:
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	var dialect gorm.Dialector
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		Logger: newLogger,
	}

	if driverDB == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			host,
			port,
			user,
			database,
			password,
			"disable",
		)

		dialect = postgres.Open(dsn)
	} else if driverDB == "mysql" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			user,
			password,
			host,
			port,
			database,
		)

		dialect = mysql.Open(dsn)
	} else if driverDB == "sqlite" {
		dialect = sqlite.Open(database)
	}

	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(20)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	//pool time
	tm := time.Minute * time.Duration(20)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(tm)

	return db, nil
}
