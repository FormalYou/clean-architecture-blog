package gorm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	)

// DSNConfig (Data Source Name) for the database connection
type DSNConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// NewDB creates a new database connection and performs auto-migration.
func NewDB(cfg DSNConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log all SQL queries
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema using GORM models
	err = db.AutoMigrate(&UserModel{}, &ArticleModel{}, &TagModel{}, &CommentModel{})
	if err != nil {
		return nil, err
	}

	return db, nil
}