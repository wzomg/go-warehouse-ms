package repository

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Name      string `mapstructure:"name"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parseTime"`
	Loc       string `mapstructure:"loc"`
}

func NewDB(cfg DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

var (
	dbOnce sync.Once
	dbInst *gorm.DB
	dbErr  error
)

func GetDB(cfg DatabaseConfig) (*gorm.DB, error) {
	dbOnce.Do(func() {
		dbInst, dbErr = NewDB(cfg)
	})
	return dbInst, dbErr
}
