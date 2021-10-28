package storage

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var store *Store

type Config struct {
	SQLiteDB     string `yaml:"sqlite_db"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Debug        bool   `yaml:"debug"`
}

type Store struct {
	db *gorm.DB
}

func (s *Store) GetDB() *gorm.DB {
	return s.db
}

func NewStore(c *Config) (*Store, error) {
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
		db.Debug()
	}
	store = &Store{db}
	return store, nil
}

func GetStore() *Store {
	return store
}
