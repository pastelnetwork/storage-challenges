package storage

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var s Store

type Config struct {
	SQLiteDB     string `mapstructure:"sqlite_db"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	Debug        bool   `mapstructure:"debug,omitempty"`
}

type Store interface {
	GetDB() *gorm.DB
}

type store struct {
	db *gorm.DB
}

func (s *store) GetDB() *gorm.DB {
	return s.db
}

func NewStore(c Config) (Store, error) {
	db, err := gorm.Open(sqlite.Open(c.SQLiteDB), &gorm.Config{SkipDefaultTransaction: true, CreateBatchSize: 200})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	maxIdleConns := 3
	maxOpenConns := 15
	maxLifeTime := 2 * time.Hour
	if c.MaxIdleConns != 0 {
		maxIdleConns = c.MaxIdleConns
	}
	if c.MaxOpenConns != 0 {
		maxOpenConns = c.MaxOpenConns
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(maxLifeTime)

	if c.Debug {
		db.Config.Logger = logger.Default.LogMode(logger.Info)
		db.Debug()
	}
	s = &store{db}

	log.Println("CONNECTED TO SQLITE DB", c)
	return s, nil
}

func GetStore() Store {
	return s
}
