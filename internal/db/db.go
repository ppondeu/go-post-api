package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ppondeu/go-post-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type sqlLogger struct {
	logger.Interface
}

func (s sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sqlString, _ := fc()
	fmt.Printf("\n===============================\n%v\n===============================\n", sqlString)
}

func ConnectDatabase(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", cfg.DB_HOST, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: sqlLogger{logger.Default.LogMode(logger.Info)},
		DryRun: false,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db
}
