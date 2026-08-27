package repository

import (
	"fmt"

	"github.com/kuonji/blog/config"
	"github.com/kuonji/blog/internal/model"
	"github.com/kuonji/blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func InitDatabase(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	logger.Info("Database connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	return nil
}

func AutoMigrate() error {
	return db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Category{},
		&model.Tag{},
		&model.Series{},
		&model.Comment{},
		&model.VisitLog{},
	)
}

func GetDB() *gorm.DB {
	return db
}

// SetTestDB 用于测试时注入内存数据库
func SetTestDB(testDB *gorm.DB) {
	db = testDB
}
